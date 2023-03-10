# to be run on the nuc only
sudo git fetch origin main    
sudo git reset --hard origin/main    
sudo systemctl daemon-reload    
sudo systemctl restart home-streaming-server.service
echo 'Successful Deploy'
