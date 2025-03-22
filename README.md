# url-shortener

## Tech Stack

File Storage: AWS S3

Reverse Proxy: Nginx

Backend: Go, Postgres

CA: Let's Encrypt (CertBot)

## Design Considerations

### Architecture
If budget wasn’t a constraint, I would have designed my architecture for high 
availability, scalability, and performance using AWS managed services.

For the frontend, I would have used Amazon S3 to store static assets, combined with 
CloudFront as a CDN to distribute content efficiently worldwide, reducing latency. 
This setup ensures faster load times and offloads traffic from backend services.

For request routing, I would have fronted my domain with either an API Gateway or an 
Application Load Balancer (ALB) to manage incoming traffic. This would allow better 
routing to different services (frontend, backend, authentication, etc.) and improve 
security and performance.

On the backend, I would have used Amazon Aurora (or RDS) for the database. Aurora 
offers automatic scaling, high availability, automated backups, and disaster 
recovery, making it a more robust solution than self-managed databases.

However, due to AWS Free Tier limitations, I opted for a more cost-effective approach 
by deploying everything on a single EC2 instance. This instance:
- Acts as a reverse proxy, handling incoming traffic.
- Hosts static frontend files directly.
- Runs multiple containers for the backend server and database.

This setup allows me to keep costs minimal while still maintaining some level of 
separation between services. However, it comes with trade-offs like limited 
scalability, single point of failure, and manual management of backups and failover.

### Short URL Collision

An alphanumeric character set typically includes:
- 26 uppercase letters (A-Z)
- 26 lowercase letters (a-z)
- 10 digits (0-9)

Total: 62 characters

The current url shortener is configured to use a 10-character alphanumerical string, 
the total number of possible unique values is:

$N = 62^{10}$

$N \approx 8.39 \times 10^{17}$

Using the Birthday Paradox Approximation, the probability of at least one collision among
𝑘 randomly chosen values is:

$P \approx 1 - e^{-\frac{k^2}{2N}}$

#### For 1 million randomly generated URLs, the probability of a collision is negligible

$P \approx 1 - e^{-\frac{(10^6)^2}{2 \times 8.39 \times 10^{17}}}$

$P \approx 0.00006$%

#### For 1 billion randomly generated URLs, the probability of a collision is significant
$P \approx 1 - e^{-\frac{(10^9)^2}{2 \times 8.39 \times 10^{17}}}$

$P \approx 44.9$%

To mitigate this risk when site traffic increases, we can:
- Increase ID length (e.g., use 12+ characters).
- Implement a database uniqueness check to prevent collisions.

## Infrastructure Setup 
You will need to set up the following in Github Secrets
- AWS_ACCESS_KEY
- AWS_SECRET_KEY
- BASIC_AUTH_PASSWORD
- DB_PASSWORD
- EC2_HOST
- EC2_USER
- JWT_SECRET
- SSH_PRIVATE_KEY

### Generate Secrets
An easy way to generate cryptographically secure random strings is to use the following command:
```bash
python -c 'import secrets; print(secrets.token_urlsafe(32))'
```

### Deploy to EC2
Use Github Actions to deploy the application to EC2 instance

### Install & Configure CertBot
SSH to EC2 instance and run the following commands
```bash
sudo dnf install -y certbot
sudo systecmtl daemon-reload
sudo systemctl enable --now certbot-renew.timer
sudo dnf install -y python3-certbot-nginx
sudo certbot certonly --nginx
```

### AWS S3
1. Create an S3 bucket
2. Disable Bucket default ACL (to allow public read)
3. Create IAM Role in AWS that allows PutObject to the S3 bucket
4. Configure the following bucket policy
```json
{
  "Version": "2012-10-17",
  "Id": "Policy1741884123207",
  "Statement": [
    {
      "Sid": "Stmt1741884121169",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::<BUCKETNAME>/*"
    },
    {
      "Sid": "Statement1",
      "Effect": "Allow",
      "Principal": {
        "AWS": "<IAM_ARN>"
      },
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::<BUCKETNAME>/*"
    }
  ]
}
```

### Configure Swap Space
AWS Free Tier instances have limited memory of 1GB. 

To prevent OOM errors, configure swap space.

SSH to EC2 instance and run the following commands
```bash
sudo dd if=/dev/zero of=/swapfile bs=128M count=16 # 2GB swap space (16*128MB)
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
sudo swapon -s # check whether swap space is enabled
sudo echo '/swapfile swap swap defaults 0 0' >> /etc/fstab
```