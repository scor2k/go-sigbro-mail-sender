#!/bin/bash -x 

cd ../../../../../ansible/ 
vim templates/nomad/sigbro_mail_sender.nomad 

ansible-playbook -Dv deploy_all_nomad_jobs.yaml --ask-vault-pass
