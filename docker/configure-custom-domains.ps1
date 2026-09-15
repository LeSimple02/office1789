# Script PowerShell pour configurer docker-mailserver

Write-Host "Configuration de docker-mailserver..." -ForegroundColor Cyan

# Copier les fichiers
docker cp .\config\postfix-virtual-domains.cf mailserver:/tmp/docker-mailserver/postfix-virtual-domains.cf
docker cp .\config\postfix-virtual-mailboxes.cf mailserver:/tmp/docker-mailserver/postfix-virtual-mailboxes.cf

Write-Host "Configuration de Postfix..." -ForegroundColor Yellow

# Verifier si la config existe
$checkResult = docker exec mailserver grep "Custom domain support" /etc/postfix/main.cf 2>&1
if ($checkResult -match "Custom domain support") {
    Write-Host "Configuration deja presente" -ForegroundColor Green
} else {
    Write-Host "Remplacement de la configuration..." -ForegroundColor Yellow
    
    # Remplacer les anciennes configurations
    docker exec mailserver sed -i 's|^virtual_mailbox_domains = /etc/postfix/vhost|virtual_mailbox_domains = regexp:/tmp/docker-mailserver/postfix-virtual-domains.cf|g' /etc/postfix/main.cf
    docker exec mailserver sed -i 's|^virtual_mailbox_maps = texthash:/etc/postfix/vmailbox|virtual_mailbox_maps = hash:/tmp/docker-mailserver/postfix-virtual-mailboxes.cf|g' /etc/postfix/main.cf
    docker exec mailserver sed -i 's|^virtual_alias_maps = texthash:/etc/postfix/virtual|virtual_alias_maps = hash:/tmp/docker-mailserver/postfix-virtual-mailboxes.cf|g' /etc/postfix/main.cf
    
    # Ajouter un commentaire
    docker exec mailserver sh -c "echo '' >> /etc/postfix/main.cf"
    docker exec mailserver sh -c "echo '# Custom domain support (Office1789)' >> /etc/postfix/main.cf"
}

Write-Host "Creation hash database..." -ForegroundColor Yellow
docker exec mailserver postmap /tmp/docker-mailserver/postfix-virtual-mailboxes.cf

Write-Host "Reload Postfix..." -ForegroundColor Yellow
docker exec mailserver postfix reload

Write-Host ""
Write-Host "Configuration terminee!" -ForegroundColor Green
Write-Host "Verifier: docker exec mailserver postconf virtual_mailbox_domains virtual_mailbox_maps virtual_alias_maps"
