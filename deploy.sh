#!/bin/bash -x 

cd ../../../../../ansible/ 
vim deploy_go_sigbro_mail_sender.yml
ansible-playbook -Dv deploy_go_sigbro_mail_sender.yml --ask-vault-pass

