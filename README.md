# CallingBoard

CallingBoard is a tool to assist ward and branch leadership in The Church of Jesus Christ of Latter-Day Saints with 
the task of reorganization.

It provides a user interface for leadership
to model a collection of calling releases and
sustainings and provides a printed report.

Data is imported from reports available on
ChurchOfJesusChrist.org within LCR. The Import tab
within the app has import links and instructions.

## USAGE
callingboard [-data <data path and html files>] [-listen :80]

By default, *callingboard* provides an HTTPS interface
on port 443. For testing locally, specify port 80 as shown above.

### Data directory
The data directory is a place where callingboard stores imported unit data. 
It also serves as the web root for the built-in http server.

#### Preparing the data directory
From the source *web* directory, copy the *html* directory and all its 
children into the data directory.

### Certificates
Once you obtain TLS certificates, you will need to hard-code the
cert path and file names in *cmd/main.go*.
secPath is the path where the certs reside.
certPath is the name of the TLS cert file, and keyPath is the name of the private key file.

There are several free sources to generate signed TLS certificates. 
I have used https://letsencrypt.org successfully, but not without much weeping
and wailing, and occasional gnashing of teeth. I really have no idea
why certificate people cannot be straightforward, but I found
every step of the way quite frustrating.

So I will provide a few bullet points to possibly save you from
the pain and suffering I endured when first exploring this realm
of the outer limits.

The following works on an Amazon EC2 instance:

1. Install yum so you can install certbot.

```
$ amazon-linux-extras install epel -y
$ yum update -y --skip-broken
$ yum install -y certbot
$ sudo yum install certbot
```

2. Clone the letsencrypt repo.

```
git clone https://github.com/letsencrypt/letsencrypt /opt/letsencrypt
```

3. Generating the certs.

    The paths in this step will vary depending on your configuration.
    
    The *-w* parameter is the web root which is your data directory plus the *html* directory.
    
    The *-d* parameter is your domain.

    This cert generation process will write some challenges to the webroot to verify that your site is really serving the domain you claim it is.

```
sudo certbot certonly --webroot -w /home/ec2-user/callingboard/html -d www.callingboard.org
```

4. The certs are only valid for 90 days. To make sure the configuration is 
valid and will endure a renew process, you can execute the following:

```
sudo certbot renew --dry-run
```
