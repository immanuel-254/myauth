package internal

import (
	"fmt"
	"strings"
	"time"
)

type Element struct {
	Name   string
	Double bool
	Attrs  map[string]string
}

func (element *Element) String(c []string) string {
	return ToString(*element, c)
}

func ToString(p Element, c []string) string {
	if !p.Double {
		var attributes []string

		if p.Attrs == nil {
			return fmt.Sprintf("<%s> </%s>", p.Name, p.Name)
		}

		for key, value := range p.Attrs {
			attributes = append(attributes, fmt.Sprintf(`%s="%s"`, key, value))
		}
		return fmt.Sprintf("<%s %s/>", p.Name, strings.Join(attributes, " "))
	}

	if p.Attrs == nil {

		if c == nil {
			return fmt.Sprintf("<%s> </%s>", p.Name, p.Name)
		}

		return fmt.Sprintf("<%s>%s</%s>", p.Name, strings.Join(c, ""), p.Name)
	}

	var attributes []string
	for key, value := range p.Attrs {
		attributes = append(attributes, fmt.Sprintf(`%s="%s"`, key, value))
	}

	if c == nil {
		return fmt.Sprintf("<%s %s> </%s>", p.Name, strings.Join(attributes, " "), p.Name)
	}

	return fmt.Sprintf("<%s %s>%s</%s>", p.Name, strings.Join(attributes, " "), strings.Join(c, ""), p.Name)
}

var (
	html = Element{
		Name:   "html",
		Double: true,
		Attrs:  map[string]string{"lang": "en"},
	}
)

func BaseHeadElement(title string) string {
	head := Element{
		Name:   "head",
		Double: true,
	}

	meta1 := Element{
		Name:   "meta",
		Double: false,
		Attrs:  map[string]string{"charset": "UTF-8"},
	}

	meta2 := Element{
		Name:   "meta",
		Double: false,
		Attrs: map[string]string{
			"name":    "viewport",
			"content": "width=device-width, initial-scale=1.0",
		},
	}

	titleEl := Element{
		Name:   "title",
		Double: true,
	}

	titleElStr := titleEl.String([]string{title})

	style := Element{
		Name:   "style",
		Double: true,
	}

	styleStr := style.String([]string{`/* TailwindCSS classes for email */
        @import url('https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css');

        /* Add inline styles for email compatibility */
        body {
            background-color: #f3f4f6;
            font-family: 'Arial', sans-serif;
        }
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
        .btn-primary {
            display: inline-block;
            padding: 12px 24px;
            background-color: #3b82f6;
            color: #ffffff;
            text-decoration: none;
            border-radius: 6px;
            font-weight: bold;
        }
        .btn-primary:hover {
            background-color: #2563eb;
        }`})

	return head.String([]string{
		meta1.String(nil),
		meta2.String(nil),
		titleElStr,
		styleStr,
	})
}

func EmailVerificationTemplate(route string) string {
	body := Element{
		Name:   "body",
		Double: true,
	}

	bodystr := body.String([]string{fmt.Sprintf(`<div class="email-container mx-auto p-6">
        <div class="text-center">
            <h1 class="text-xl font-bold text-gray-800">Verify Your Email Address</h1>
        </div>
        <div class="mt-6">
            <p class="text-gray-600 text-sm">
                Thank you for signing up! Please verify your email address to complete your registration.
            </p>
        </div>
        <div class="mt-6 text-center">
            <a href="%s" class="btn-primary">Verify Email</a>
        </div>
        <div class="mt-6 text-sm text-gray-500">
            <p>If you didn’t sign up for this account, you can safely ignore this email.</p>
        </div>
        <div class="mt-6 text-center text-xs text-gray-400">
            <p>&copy; %v Your Company. All rights reserved.</p>
        </div>
    </div>`, route, time.Now().Year())})

	return html.String([]string{BaseHeadElement("Email Verification"), bodystr})
}

func ChangeEmailVerificationTemplate(route string) string {
	body := Element{
		Name:   "body",
		Double: true,
	}

	bodystr := body.String([]string{fmt.Sprintf(`<div class="email-container mx-auto p-6">
        <div class="text-center">
            <h1 class="text-xl font-bold text-gray-800">Change Your Email Address</h1>
        </div>
        <div class="mt-6">
            <p class="text-gray-600 text-sm">
                You are receiving this email because you have requested to change your email address attributed to your account. 
				Click on the link below to verify that you wish to make this changes
            </p>
        </div>
        <div class="mt-6 text-center">
            <a href="%s" class="btn-primary">Change Email</a>
        </div>
        <div class="mt-6 text-sm text-gray-500">
            <p>If you did not request to change your email address, you can safely ignore this email.</p>
        </div>
        <div class="mt-6 text-center text-xs text-gray-400">
            <p>&copy; %v Your Company. All rights reserved.</p>
        </div>
    </div>`, route, time.Now().Year())})

	return html.String([]string{BaseHeadElement("Change Email Verification"), bodystr})
}

func ChangePasswordVerificationTemplate(route string) string {
	body := Element{
		Name:   "body",
		Double: true,
	}

	bodystr := body.String([]string{fmt.Sprintf(`<div class="email-container mx-auto p-6">
        <div class="text-center">
            <h1 class="text-xl font-bold text-gray-800">Change Your Password</h1>
        </div>
        <div class="mt-6">
            <p class="text-gray-600 text-sm">
                You are receiving this email because you have requested to change your Password. 
				Click on the link below to verify that you wish to make this changes
            </p>
        </div>
        <div class="mt-6 text-center">
            <a href="%s" class="btn-primary">Change Password</a>
        </div>
        <div class="mt-6 text-sm text-gray-500">
            <p>If you did not request to change your password, you can safely ignore this email.</p>
        </div>
        <div class="mt-6 text-center text-xs text-gray-400">
            <p>&copy; %v Your Company. All rights reserved.</p>
        </div>
    </div>`, route, time.Now().Year())})

	return html.String([]string{BaseHeadElement("Change Password Verification"), bodystr})
}

func ResetPasswordVerificationTemplate(route string) string {
	body := Element{
		Name:   "body",
		Double: true,
	}

	bodystr := body.String([]string{fmt.Sprintf(`<div class="email-container mx-auto p-6">
        <div class="text-center">
            <h1 class="text-xl font-bold text-gray-800">Reset Your Password</h1>
        </div>
        <div class="mt-6">
            <p class="text-gray-600 text-sm">
                You are receiving this email because you have requested to Reset your Password. 
				Click on the link below to verify that you wish to make this changes
            </p>
        </div>
        <div class="mt-6 text-center">
            <a href="%s" class="btn-primary">Reset Password</a>
        </div>
        <div class="mt-6 text-sm text-gray-500">
            <p>If you did not request to reset your password, you can safely ignore this email.</p>
        </div>
        <div class="mt-6 text-center text-xs text-gray-400">
            <p>&copy; %v Your Company. All rights reserved.</p>
        </div>
    </div>`, route, time.Now().Year())})

	return html.String([]string{BaseHeadElement("Reset Password Verification"), bodystr})
}

func DeleteUserVerificationTemplate(route string) string {
	body := Element{
		Name:   "body",
		Double: true,
	}

	bodystr := body.String([]string{fmt.Sprintf(`<div class="email-container mx-auto p-6">
        <div class="text-center">
            <h1 class="text-xl font-bold text-gray-800">Delete User Account</h1>
        </div>
        <div class="mt-6">
            <p class="text-gray-600 text-sm">
                You are receiving this email because you have requested to delete your user account. 
				Click on the link below to verify that you wish to make this changes
            </p>
        </div>
        <div class="mt-6 text-center">
            <a href="%s" class="btn-primary">Delete User Account</a>
        </div>
        <div class="mt-6 text-sm text-gray-500">
            <p>If you did not request to delete user account, you can safely ignore this email.</p>
        </div>
        <div class="mt-6 text-center text-xs text-gray-400">
            <p>&copy; %v Your Company. All rights reserved.</p>
        </div>
    </div>`, route, time.Now().Year())})

	return html.String([]string{BaseHeadElement("Reset Password Verification"), bodystr})
}
