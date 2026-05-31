import smtplib
from email.mime.text import MIMEText

GMAIL_USER = "noreplycastonion@gmail.com"
GMAIL_APP_PASSWORD = "pkdg jccu qqaf nzqp"

SMTP_SERVER = "smtp.gmail.com"
SMTP_PORT = 465

EMAIL_TEMPLATE = """\
Hello {name},

We are pleased to inform you that your application to broadcast on cast.onion has been approved.

Your station "{station_name}" is now active on the network.

---

Your credentials are below. Please save them securely — they will not be sent again.

Station Key (use this to authenticate your stream to the broadcast server):

  {station_key}

Access Token (use this to log into the station owner dashboard):

  {access_token}

---

Getting started:

1. Download the cast.onion desktop client
2. Open the Broadcast tab and enter your Access Token
3. Enter your Station Key and select your microphone
4. Click "microphone" to go live

If you have any questions, reply to this email.

---

cast.onion
private internet radio
"""

def send_approval_email(to_email, name, station_name, station_key, access_token):
    body = EMAIL_TEMPLATE.format(
        name=name,
        station_name=station_name,
        station_key=station_key,
        access_token=access_token
    )

    msg = MIMEText(body, "plain")
    msg["From"] = GMAIL_USER
    msg["To"] = to_email
    msg["Subject"] = "Your cast.onion application has been approved"

    with smtplib.SMTP_SSL(SMTP_SERVER, SMTP_PORT) as server:
        server.login(GMAIL_USER, GMAIL_APP_PASSWORD)
        server.send_message(msg)

    print(f"Email sent to {to_email}")

if __name__ == "__main__":
    send_approval_email(
        to_email="marioplays254@gmail.com",
        name="Luis Vega",
        station_name="lofi-fm",
        station_key="station_abc123",
        access_token="token_xyz456"
    )