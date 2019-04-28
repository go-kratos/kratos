package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/macross/conf"
	"go-common/app/admin/main/macross/model/mail"
	"go-common/app/admin/main/macross/tools"
	"go-common/library/log"
	"io"
	"os"
	"path/filepath"
	"strings"

	gomail "gopkg.in/gomail.v2"
)

// SendMail send mail
func (s *Service) SendMail(c context.Context, m *mail.Mail, attach *mail.Attach) (err error) {
	var (
		toUsers  []string
		ccUsers  []string
		bccUsers []string
		msg      = gomail.NewMessage()
	)

	msg.SetAddressHeader("From", conf.Conf.Property.Mail.Address, conf.Conf.Property.Mail.Name) // 发件人
	for _, ads := range m.ToAddresses {
		toUsers = append(toUsers, msg.FormatAddress(ads.Address, ads.Name))
	}

	for _, ads := range m.CcAddresses {
		ccUsers = append(ccUsers, msg.FormatAddress(ads.Address, ads.Name))
	}

	for _, ads := range m.BccAddresses {
		bccUsers = append(bccUsers, msg.FormatAddress(ads.Address, ads.Name))
	}

	msg.SetHeader("To", toUsers...)
	msg.SetHeader("Subject", m.Subject) // 主题

	if len(ccUsers) > 0 {
		msg.SetHeader("Cc", ccUsers...)
	}
	if len(bccUsers) > 0 {
		msg.SetHeader("Bcc", bccUsers...)
	}

	if m.Type == mail.TypeTextHTML {
		msg.SetBody("text/html", m.Body)
	} else {
		msg.SetBody("text/plain", m.Body)
	}

	// 附件处理
	if attach != nil {
		tmpSavePath := filepath.Join(os.TempDir(), "mail_tmp")
		err = os.MkdirAll(tmpSavePath, 0755)
		if err != nil {
			log.Error("os.MkdirAll error(%v)", err)
			return
		}
		destFilePath := filepath.Join(tmpSavePath, attach.Name)
		destFile, cErr := os.Create(destFilePath)
		if cErr != nil {
			log.Error("os.Create(%s) error(%v)", destFilePath, cErr)
			return cErr
		}
		defer os.RemoveAll(tmpSavePath)
		io.Copy(destFile, attach.File)

		// 如果 zip 文件需要解压以后放在邮件附件中
		if attach.ShouldUnzip && strings.HasSuffix(attach.Name, ".zip") {
			unzipFilePath := filepath.Join(tmpSavePath, "unzip")
			err = os.MkdirAll(tmpSavePath, 0755)
			if err != nil {
				log.Error("os.MkdirAll error(%v)", err)
				return
			}
			err = tools.Unzip(destFilePath, unzipFilePath)
			if err != nil {
				log.Error("unzip(%s, %s) error(%v)", destFilePath, unzipFilePath, err)
				return
			}
			err = filepath.Walk(unzipFilePath, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					log.Error("filepath.Walk error(%v)", err)
					return err
				}
				if f == nil {
					errMsg := "found no file"
					err = fmt.Errorf(errMsg)
					log.Error(errMsg)
					return err
				}
				if f.IsDir() {
					return nil
				}
				msg.Attach(path)
				return err
			})
		} else {
			msg.Attach(destFilePath)
		}
	}

	d := gomail.NewDialer(
		conf.Conf.Property.Mail.Host,
		conf.Conf.Property.Mail.Port,
		conf.Conf.Property.Mail.Address,
		conf.Conf.Property.Mail.Pwd,
	)
	if err = d.DialAndSend(msg); err != nil {
		log.Error("Send mail Fail(%v) diff(%s)", msg, err)
		return
	}

	return
}
