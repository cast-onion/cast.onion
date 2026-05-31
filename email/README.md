# Email Sender

This Rust folder is where the function for sending emails for whether the station applicant's application got denied or approved, or if their station keys got suspended or revoked.

## Email

This will always send under noreplycastonion@gmail.com and no other emails.

If you are contributing to this section of the service and want to try with your gmail, follow these steps:

### Script

The script to test it out is in scripts/send_email.py where you can change the email, app password, and the email you want to send to.

To get an app password, make sure you have 2FA enabled on your Google account.

Then go to here, https://myaccount.google.com/apppasswords and put a name for it.

After putting in the your app password's name, you will be given your special app password.

Do not give this app password to anybody or they could potentially attempt to access your Google account.

Once you have your app password, please change:

- `GMAIL_USER = "example@gmail.com"`
- `GMAIL_APP_PASSWORD = "your app password"`

Once changed, you may run

```bash
python scripts/send_email.py
```

or

```bash
python3 scripts/send_email.py
```

No need to install anything as the modules imported are pre-added with Pip/Python installation.