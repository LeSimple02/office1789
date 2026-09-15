#!/bin/bash
# Script to configure docker-mailserver for custom domains support

echo "Configuring docker-mailserver for custom domains..."

# Copy virtual domains and mailboxes to mailserver container
echo "Copying virtual domains configuration..."
docker cp ./config/postfix-virtual-domains.cf mailserver:/tmp/docker-mailserver/postfix-virtual-domains.cf

echo "Copying virtual mailboxes configuration..."
docker cp ./config/postfix-virtual-mailboxes.cf mailserver:/tmp/docker-mailserver/postfix-virtual-mailboxes.cf

# Check if configuration already exists
if docker exec mailserver grep -q "Custom domain support" /etc/postfix/main.cf 2>/dev/null; then
    echo "Configuration already present"
else
    echo "Replacing Postfix configuration..."
    
    # Replace old configurations with new ones
    docker exec mailserver sed -i 's|^virtual_mailbox_domains = /etc/postfix/vhost|virtual_mailbox_domains = regexp:/tmp/docker-mailserver/postfix-virtual-domains.cf|g' /etc/postfix/main.cf
    docker exec mailserver sed -i 's|^virtual_mailbox_maps = texthash:/etc/postfix/vmailbox|virtual_mailbox_maps = hash:/tmp/docker-mailserver/postfix-virtual-mailboxes.cf|g' /etc/postfix/main.cf
    docker exec mailserver sed -i 's|^virtual_alias_maps = texthash:/etc/postfix/virtual|virtual_alias_maps = hash:/tmp/docker-mailserver/postfix-virtual-mailboxes.cf|g' /etc/postfix/main.cf
    
    # Add comment
    docker exec mailserver sh -c "echo '' >> /etc/postfix/main.cf"
    docker exec mailserver sh -c "echo '# Custom domain support (Office1789)' >> /etc/postfix/main.cf"
fi

# Create Postfix hash database for virtual mailboxes
echo "Creating Postfix hash database..."
docker exec mailserver postmap /tmp/docker-mailserver/postfix-virtual-mailboxes.cf

# Reload Postfix
echo "Reloading Postfix configuration..."
docker exec mailserver postfix reload

echo ""
echo "Configuration completed successfully!"
echo "Verify: docker exec mailserver postconf virtual_mailbox_domains virtual_mailbox_maps virtual_alias_maps"
