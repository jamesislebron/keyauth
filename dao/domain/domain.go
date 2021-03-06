package domain

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/defineiot/keyauth/internal/exception"
)

const (
	Personal   Type = iota // Personal 个人域
	Enterprise             //  Enterprise 企业域
	Paterner               // 合作伙伴域
)

// Type 域类型
type Type int

// Domain a tenant container, example an company or organization.
type Domain struct {
	ID             string `json:"id"`                        // 域ID
	Type           Type   `json:"type,omitempty"`            // 域类型: Personal: 个人, Enterprise: 企业, Paterner: 合作伙伴伙伴
	Name           string `json:"name,omitempty"`            // 公司或者组织名称
	DisplayName    string `json:"display_name,omitempty"`    // 全称
	LogoPath       string `json:"logo_path,omitempty"`       // 公司LOGO图片的URL
	Description    string `json:"description,omitempty"`     // 描述
	Enabled        bool   `json:"enabled,omitempty"`         // 域状态, 是否需要冻结该域, 冻结时, 该域下面所有用户禁止登录
	Size           string `json:"size,omitempty"`            // 规模: 50人以下, 50~100, ...
	Location       string `json:"location,omitempty"`        // 位置: 指城市, 比如 中国,四川,成都
	Address        string `json:"address,omitempty"`         // 地址: 比如环球中心 10F 1034
	Industry       string `json:"industry,omitempty"`        // 所属行业: 比如, 互联网
	Fax            string `json:"fax,omitempty"`             // 传真:
	Phone          string `json:"phone,omitempty"`           // 固话:
	ContactsName   string `json:"contacts_name,omitempty"`   // 联系人姓名
	ContactsTitle  string `json:"contacts_title,omitempty"`  // 联系人职位
	ContactsMobile string `json:"contacts_mobile,omitempty"` // 联系人电话
	ContactsEmail  string `json:"contacts_email,omitempty"`  // 联系人邮箱
	CreateAt       int64  `json:"create_at,omitempty"`       // 创建时间
	UpdateAt       int64  `json:"update_at,omitempty"`       // 更新时间
	Owner          string `json:"owner,omitempty"`           // 创建者ID
}

func (d *Domain) String() string {
	str, err := json.Marshal(d)
	if err != nil {
		log.Printf("E! marshal domain to string error: %s", err)
		return fmt.Sprintf("ID: %s, Name: %s", d.ID, d.Name)
	}

	return string(str)
}

// Validate 创建时校验参数合法性
func (d *Domain) Validate() error {
	if d.Name == "" {
		return exception.NewBadRequest("domain's name is required!")
	}

	if len(d.Name) > 128 {
		return exception.NewBadRequest("domain's name is too long,  max length is 128")
	}

	return nil
}

// Store is an domain service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader for read data from store
type Reader interface {
	GetDomainByID(domainID string) (*Domain, error)
	GetDomainByName(name string) (*Domain, error)
	CheckDomainIsExistByID(domainID string) (bool, error)
	CheckDomainIsExistByName(domainName string) (bool, error)
	ListDomain(pageNumber, pageSize int64) (domains []*Domain, totalPage int64, err error)
	ListUserThirdDomains(userID string) ([]*Domain, error)
}

// Writer for write data to store
type Writer interface {
	CreateDomain(d *Domain) error
	UpdateDomain(id, name, description string) (*Domain, error)
	DeleteDomainByID(id string) error
	DeleteDomainByName(name string) error
}
