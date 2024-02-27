package client

import (
	"bytes"
	"debug/elf"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	URL string
}

type PingRequest struct{}
type PingReply struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text    string `xml:",chardata"`
				Boolean string `xml:"boolean"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

func (c *Client) Ping() (*PingReply, error) {
	req := []byte("<methodCall><methodName>d2d.ping</methodName><params></params></methodCall>")

	resp, err := http.Post(c.URL, "text/xml", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to send ping request: %w", err)
	}
	defer resp.Body.Close()

	reply := &PingReply{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(body, reply)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ping response: %w", err)
	}

	return reply, nil
}

type FunctionHeadersRequest struct{}

type FunctionHeader struct {
	Name  string
	Size  int
	Value int
}

type FunctionHeadersReply struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text   string `xml:",chardata"`
							Struct struct {
								Text   string `xml:",chardata"`
								Member []struct {
									Text  string `xml:",chardata"`
									Name  string `xml:"name"`
									Value struct {
										Text string `xml:",chardata"`
										I4   string `xml:"i4"`
									} `xml:"value"`
								} `xml:"member"`
							} `xml:"struct"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

func (c *Client) FunctionHeaders() ([]*FunctionHeader, error) {
	req := []byte("<methodCall><methodName>d2d.function_headers</methodName><params></params></methodCall>")

	resp, err := http.Post(c.URL, "text/xml", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to send function headers request: %w", err)
	}
	defer resp.Body.Close()

	reply := &FunctionHeadersReply{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read function vars response body")
	}

	err = xml.Unmarshal(body, reply)
	if err != nil {
		panic(err)
	}

	result := []*FunctionHeader{}

	for _, v := range reply.Params.Param.Value.Struct.Member {
		name, _ := strings.CutPrefix(v.Name, "0x")
		value, err2 := strconv.ParseUint(name, 16, 32)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse function address %s: %w", v.Name, err2)
		}

		fh := &FunctionHeader{
			Value: int(value),
		}

		for _, p := range v.Value.Struct.Member {
			switch p.Name {
			case "name":
				fh.Name = p.Value.Text
			case "size":
				size, err3 := strconv.Atoi(p.Value.I4)
				if err3 != nil {
					return nil, fmt.Errorf("cannot parse function size %s: %w", p.Value.I4, err3)
				}

				fh.Size = size

			}
		}

		result = append(result, fh)
	}

	return result, nil
}

type GlobalVarsReply struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text   string `xml:",chardata"`
							Struct struct {
								Text   string `xml:",chardata"`
								Member struct {
									Text  string `xml:",chardata"`
									Name  string `xml:"name"`
									Value string `xml:"value"`
								} `xml:"member"`
							} `xml:"struct"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type GlobalVar struct {
	Name  string
	Value int
}

func (c *Client) GlobalVars() ([]*GlobalVar, error) {
	req := []byte("<methodCall><methodName>d2d.global_vars</methodName><params></params></methodCall>")

	resp, err := http.Post(c.URL, "text/xml", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to send global vars request: %w", err)
	}
	defer resp.Body.Close()

	reply := &GlobalVarsReply{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read global vars response body")
	}

	err = xml.Unmarshal(body, reply)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal global vars: %w", err)
	}

	result := []*GlobalVar{}

	for _, v := range reply.Params.Param.Value.Struct.Member {
		vv, _ := strings.CutPrefix(v.Name, "0x")
		value, err2 := strconv.ParseUint(vv, 16, 32)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", v.Name, err2)
		}

		fh := &GlobalVar{
			Value: int(value),
			Name:  v.Value.Struct.Member.Value,
		}

		result = append(result, fh)
	}

	return result, nil
}

type ImageBaseReply struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

func (c *Client) GetImageBase() (int, error) {
	req := []byte("<methodCall><methodName>d2d.getImageBase</methodName><params></params></methodCall>")

	resp, err := http.Post(c.URL, "text/xml", bytes.NewBuffer(req))
	if err != nil {
		return 0, fmt.Errorf("failed to send ping request: %w", err)
	}
	defer resp.Body.Close()

	reply := &ImageBaseReply{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = xml.Unmarshal(body, reply)
	if err != nil {
		return 0, fmt.Errorf("failed to decode ping response: %w", err)
	}

	v, _ := strings.CutPrefix(reply.Params.Param.Value, "0x")
	base, err := strconv.ParseInt(v, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s into int: %w", reply.Params.Param.Value, err)
	}

	return int(base), nil
}

type ElfInfoResp struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text    string `xml:",chardata"`
							I4      int    `xml:"i4"`
							Boolean bool   `xml:"boolean"`
							Value   string `xml:"value"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type ElfInfo struct {
	Machine     elf.Machine
	ImageBase   int
	Error       error
	Flags       int
	IsBigEndian bool
	Is32Bit     bool
	Name        string
}

func HexToInt(hex string) (int, error) {
	v, _ := strings.CutPrefix(hex, "0x")
	i, err := strconv.ParseInt(v, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s into int: %w", hex, err)
	}

	return int(i), nil
}

func (c *Client) ElfInfo() (*ElfInfo, error) {
	req := []byte("<methodCall><methodName>d2d.elf_info</methodName><params></params></methodCall>")

	resp, err := http.Post(c.URL, "text/xml", bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("failed to send ping request: %w", err)
	}
	defer resp.Body.Close()

	reply := &ElfInfoResp{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(body, reply)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ping response: %w", err)
	}

	elfInfo := &ElfInfo{
		Machine:     0,
		ImageBase:   0,
		Error:       nil,
		Flags:       0,
		IsBigEndian: false,
		Is32Bit:     false,
		Name:        "",
	}

	for _, m := range reply.Params.Param.Value.Struct.Member {
		switch m.Name {
		case "error":
			if m.Value.Text != "" {
				return nil, errors.New(m.Value.Text)
			}
		case "machine":
			elfInfo.Machine = elf.Machine(m.Value.I4)
		case "is_big_endian":
			elfInfo.IsBigEndian = m.Value.Boolean
		case "flags":
			elfInfo.Flags, err = HexToInt(m.Value.Text)
			if err != nil {
				return nil, fmt.Errorf("cannot parse elf flags %s: %w", m.Value.Text, err)
			}
		case "image_base":
			elfInfo.ImageBase, err = HexToInt(m.Value.Text)
			if err != nil {
				return nil, fmt.Errorf("cannot parse image base %s: %w", m.Value.Text, err)
			}
		case "is_32_bit":
			elfInfo.Is32Bit = m.Value.Boolean
		case "name":
			elfInfo.Name = m.Value.Text
		default:
			return nil, fmt.Errorf("unexpected field %s", m.Name)
		}
	}

	return elfInfo, nil
}
