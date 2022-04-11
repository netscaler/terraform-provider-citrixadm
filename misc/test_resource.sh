set -x
echo "Lets start"
read -p "Press enter to INIT"
terraform init
read -p "Press enter to CREATE"
terraform apply -auto-approve
read -p "Press enter to check IDEMPOTENCY"
terraform apply -auto-approve
read -p "Press enter to UPDATE"
terraform apply -auto-approve
read -p "Press enter to DESTROY"
terraform destroy -auto-approve
echo "The end"