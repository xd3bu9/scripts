
# get a list network interfaces on the host device
interfaces=$(cat /proc/net/route|grep -vi iface|awk -F" " '{ print $1 }'|uniq -u)

# Extract mac addresses from the interfaces present
echo $interfaces|while read -r line;do echo "$line::$(cat /sys/class/net/$line/address)";done
