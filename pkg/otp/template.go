package otp

import (
	"github.com/authgear/authgear-server/pkg/auth/config"
	"github.com/authgear/authgear-server/pkg/auth/dependency/identity/loginid"
	"github.com/authgear/authgear-server/pkg/template"
)

type MessageOrigin string

const (
	MessageOriginSignup   MessageOrigin = "signup"
	MessageOriginLogin    MessageOrigin = "login"
	MessageOriginSettings MessageOrigin = "settings"
)

type OOBOperationType string

const (
	OOBOperationTypePrimaryAuth   OOBOperationType = "primary_auth"
	OOBOperationTypeSecondaryAuth OOBOperationType = "secondary_auth"
	OOBOperationTypeVerify        OOBOperationType = "verify"
)

type MessageTemplateContext struct {
	AppName              string
	Email                string
	Phone                string
	LoginID              *loginid.LoginID
	Code                 string
	Host                 string
	Origin               MessageOrigin
	Operation            OOBOperationType
	StaticAssetURLPrefix string
}

const (
	TemplateItemTypeOTPMessageSMSTXT    config.TemplateItemType = "otp_message_sms.txt"
	TemplateItemTypeOTPMessageEmailTXT  config.TemplateItemType = "otp_message_email.txt"
	TemplateItemTypeOTPMessageEmailHTML config.TemplateItemType = "otp_message_email.html"
)

var TemplateOTPMessageSMSTXT = template.Spec{
	Type: TemplateItemTypeOTPMessageSMSTXT,
	Default: `{{ .Code }} is your {{ .AppName }} verification code.

Please ignore if you didn't sign in or sign up.

@{{ .Host }} #{{ .Code }}
`,
}

var TemplateOTPMessageEmailTXT = template.Spec{
	Type: TemplateItemTypeOTPMessageEmailTXT,
	Default: `Verify your email on {{ .AppName }}

You have selected {{ .Email }} for verification. Please use the following code to complete to the verification.

{{ .Code }}

If you didn't sign in or sign up please ignore this email.
`,
}

var TemplateOTPMessageEmailHTML = template.Spec{
	Type:   TemplateItemTypeOTPMessageEmailHTML,
	IsHTML: true,
	Default: `
<!-- FILE: ./templates/otp_message_email.mjml -->
<!doctype html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
  <title>
  </title>
  <!--[if !mso]><!-- -->
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <!--<![endif]-->
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <style type="text/css">
    #outlook a {
      padding: 0;
    }

    body {
      margin: 0;
      padding: 0;
      -webkit-text-size-adjust: 100%;
      -ms-text-size-adjust: 100%;
    }

    table,
    td {
      border-collapse: collapse;
      mso-table-lspace: 0pt;
      mso-table-rspace: 0pt;
    }

    img {
      border: 0;
      height: auto;
      line-height: 100%;
      outline: none;
      text-decoration: none;
      -ms-interpolation-mode: bicubic;
    }

    p {
      display: block;
      margin: 13px 0;
    }
  </style>
  <!--[if mso]>
        <xml>
        <o:OfficeDocumentSettings>
          <o:AllowPNG/>
          <o:PixelsPerInch>96</o:PixelsPerInch>
        </o:OfficeDocumentSettings>
        </xml>
        <![endif]-->
  <!--[if lte mso 11]>
        <style type="text/css">
          .mj-outlook-group-fix { width:100% !important; }
        </style>
        <![endif]-->
  <style type="text/css">
    @media only screen and (min-width:480px) {
      .mj-column-per-100 {
        width: 100% !important;
        max-width: 100%;
      }

      .mj-column-px-250 {
        width: 250px !important;
        max-width: 250px;
      }
    }
  </style>
  <style type="text/css">
    @media only screen and (max-width:480px) {
      table.mj-full-width-mobile {
        width: 100% !important;
      }

      td.mj-full-width-mobile {
        width: auto !important;
      }
    }
  </style>
</head>

<body>
  <div style="">
    <!--[if mso | IE]>
      <table
         align="center" border="0" cellpadding="0" cellspacing="0" class="" style="width:600px;" width="600"
      >
        <tr>
          <td style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;">
      <![endif]-->
    <div style="margin:0px auto;max-width:600px;">
      <table align="center" border="0" cellpadding="0" cellspacing="0" role="presentation" style="width:100%;">
        <tbody>
          <tr>
            <td style="direction:ltr;font-size:0px;padding:20px 0;text-align:center;">
              <!--[if mso | IE]>
                  <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                
        <tr>
      
            <td
               class="" style="vertical-align:top;width:600px;"
            >
          <![endif]-->
              <div class="mj-column-per-100 mj-outlook-group-fix" style="font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;">
                <table border="0" cellpadding="0" cellspacing="0" role="presentation" style="vertical-align:top;" width="100%">
                  <tr>
                    <td align="center" style="font-size:0px;padding:20px;word-break:break-word;">
                      <div style="font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji;font-size:24px;font-weight:bold;line-height:1;text-align:center;color:#000000;">Verify your email on {{ .AppName }}</div>
                    </td>
                  </tr>
                  <tr>
                    <td align="center" style="font-size:0px;padding:20px;word-break:break-word;">
                      <div style="font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji;font-size:16px;line-height:1;text-align:center;color:#000000;">You have selected {{ .Email }} for verification. Please use the following code to complete the verification.</div>
                    </td>
                  </tr>
                </table>
              </div>
              <!--[if mso | IE]>
            </td>
          
        </tr>
      
                  </table>
                <![endif]-->
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <!--[if mso | IE]>
          </td>
        </tr>
      </table>
      
      <table
         align="center" border="0" cellpadding="0" cellspacing="0" class="" style="width:600px;" width="600"
      >
        <tr>
          <td style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;">
      <![endif]-->
    <div style="margin:0px auto;max-width:600px;">
      <table align="center" border="0" cellpadding="0" cellspacing="0" role="presentation" style="width:100%;">
        <tbody>
          <tr>
            <td style="direction:ltr;font-size:0px;padding:20px 0;text-align:center;">
              <!--[if mso | IE]>
                  <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                
        <tr>
      
            <td
               class="" style="vertical-align:top;width:250px;"
            >
          <![endif]-->
              <div class="mj-column-px-250 mj-outlook-group-fix" style="font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;">
                <table border="0" cellpadding="0" cellspacing="0" role="presentation" style="background-color:#f1f4f5;vertical-align:top;" width="100%">
                  <tr>
                    <td align="center" style="font-size:0px;padding:24px 24px 24px 40px;word-break:break-word;">
                      <div style="font-family:monospace;font-size:36px;font-weight:heavy;letter-spacing:16px;line-height:1;text-align:center;color:#000000;">{{ .Code }}</div>
                    </td>
                  </tr>
                </table>
              </div>
              <!--[if mso | IE]>
            </td>
          
        </tr>
      
                  </table>
                <![endif]-->
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <!--[if mso | IE]>
          </td>
        </tr>
      </table>
      
      <table
         align="center" border="0" cellpadding="0" cellspacing="0" class="" style="width:600px;" width="600"
      >
        <tr>
          <td style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;">
      <![endif]-->
    <div style="margin:0px auto;max-width:600px;">
      <table align="center" border="0" cellpadding="0" cellspacing="0" role="presentation" style="width:100%;">
        <tbody>
          <tr>
            <td style="direction:ltr;font-size:0px;padding:20px 0;text-align:center;">
              <!--[if mso | IE]>
                  <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                
        <tr>
      
            <td
               class="" style="vertical-align:top;width:600px;"
            >
          <![endif]-->
              <div class="mj-column-per-100 mj-outlook-group-fix" style="font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;">
                <table border="0" cellpadding="0" cellspacing="0" role="presentation" style="vertical-align:top;" width="100%">
                  <tr>
                    <td align="center" style="font-size:0px;padding:20px;word-break:break-word;">
                      <div style="font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji;font-size:12px;font-weight:light;line-height:1;text-align:center;color:#000000;">If you didn't sign in or sign up please ignore this email.</div>
                    </td>
                  </tr>
                  <tr>
                    <td align="center" style="font-size:0px;padding:60px;word-break:break-word;">
                      <table border="0" cellpadding="0" cellspacing="0" role="presentation" style="border-collapse:collapse;border-spacing:0px;">
                        <tbody>
                          <tr>
                            <td style="width:65px;">
                              <img height="15" src="{{ .StaticAssetURLPrefix }}/authui/image/ic_footer_authgear.png" style="border:0;display:block;outline:none;text-decoration:none;height:15px;width:100%;font-size:13px;" width="65" />
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                </table>
              </div>
              <!--[if mso | IE]>
            </td>
          
        </tr>
      
                  </table>
                <![endif]-->
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <!--[if mso | IE]>
          </td>
        </tr>
      </table>
      <![endif]-->
  </div>
</body>

</html>
`,
}
