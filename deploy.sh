# to be run on the nuc only
git fetch origin main    
git reset --hard origin/main    
systemctl daemon-reload    
systemctl restart home-streaming-server.service
echo 'Successful Deploy'
