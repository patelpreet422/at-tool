package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"syscall"

	"github.com/fatih/color"
	"github.com/patelpreet422/at-tool/cookiejar"
	"github.com/patelpreet422/at-tool/util"
	"golang.org/x/crypto/ssh/terminal"
)

// ErrorNotLogged not logged in
const ErrorNotLogged = "Not logged in"

// findHandle if logged return (handle, nil), else return ("", ErrorNotLogged)
func findUsername(body []byte) (string, error) {
	reg := regexp.MustCompile(`userScreenName = "(.+?)"`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New(ErrorNotLogged)
	}

	return string(tmp[1]), nil
}

func findCsrf(body []byte) (string, error) {
	reg := regexp.MustCompile(`csrfToken = "(.+?)"`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("Cannot find csrf")
	}

	return string(tmp[1]), nil
}

// Login codeforces with handler and password
func (c *Client) Login() (err error) {
	color.Cyan("Login %v...\n", c.Username)

	password, err := c.DecryptPassword()
	if err != nil {
		return
	}

	jar, _ := cookiejar.New(nil)
	c.client.Jar = jar
	body, err := util.GetBody(c.client, c.host+"/login")
	if err != nil {
		return
	}

	csrf, err := findCsrf(body)
	if err != nil {
		return
	}

	body, err = util.PostBody(c.client, c.host+"/login", url.Values{
		"username":   {c.Username},
		"password":   {password},
		"csrf_token": {csrf},
	})
	if err != nil {
		return
	}

	username, err := findUsername(body)
	if err != nil {
		return
	}

	c.Username = username
	c.Jar = jar
	color.Green("Succeed!!")
	color.Green("Welcome %v~", username)
	return c.save()
}

func createHash(key string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func encrypt(handle, password string) (ret string, err error) {
	block, err := aes.NewCipher(createHash("glhf" + handle + "233"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	text := gcm.Seal(nonce, nonce, []byte(password), nil)
	ret = hex.EncodeToString(text)
	return
}

func decrypt(handle, password string) (ret string, err error) {
	data, err := hex.DecodeString(password)
	if err != nil {
		err = errors.New("Cannot decode the password")
		return
	}
	block, err := aes.NewCipher(createHash("glhf" + handle + "233"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	nonce, text := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return
	}
	ret = string(plain)
	return
}

// DecryptPassword get real password
func (c *Client) DecryptPassword() (string, error) {
	if len(c.Password) == 0 || len(c.Username) == 0 {
		return "", errors.New("You have to configure your username and password by `at config`")
	}
	return decrypt(c.Username, c.Password)
}

// ConfigLogin configure handle and password
func (c *Client) ConfigLogin() (err error) {
	if c.Username != "" {
		color.Green("Current user: %v", c.Username)
	}
	color.Cyan("Configure username and password")
	color.Cyan("Note: The password is invisible, just type it correctly.")

	fmt.Printf("Username: ")
	username := util.ScanlineTrim()

	password := ""
	if terminal.IsTerminal(int(syscall.Stdin)) {
		fmt.Printf("Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println()
			if err.Error() == "EOF" {
				fmt.Println("Interrupted.")
				return nil
			}
			return err
		}
		password = string(bytePassword)
		fmt.Println()
	} else {
		color.Red("Your terminal does not support the hidden password.")
		fmt.Printf("password: ")
		password = util.Scanline()
	}

	c.Username = username
	c.Password, err = encrypt(username, password)
	if err != nil {
		return
	}
	return c.Login()
}
