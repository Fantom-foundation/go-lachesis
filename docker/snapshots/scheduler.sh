#!/bin/bash \n\

mv /docker-entrypoint.d/index.html /usr/share/nginx/html/index.html
mkdir -p /snapshots/uploads && touch /snapshots/latest.snapshot
while [ -f /docker-entrypoint.d/mainnet.toml ]
do
    previous_epoch=\$(ls -1 /snapshots/*.epoch 2>/dev/null 1>&2 && echo \$(basename \$(ls -1 /snapshots/*.epoch | sort -r | head -n 1)) || echo 0.epoch)
    previous_epoch=\${previous_epoch%%.epoch}
    lachesis --nousb --config /docker-entrypoint.d/mainnet.toml &
    pid=\$!
    sleep 333
    while [ \$(( \$(date +%%s) - \$(stat -c %%Y /snapshots/latest.snapshot) )) -lt 4321 ]
    do
        sleep 6
    done
    current_epoch=\$(lachesis attach --exec='ftm.currentEpoch()')
    kill \$pid && wait \$pid
    lachesis export /snapshots/latest.snapshot \$last_epoch \$current_epoch
    unlink /snapshots/\$previous_epoch.epoch
    ln -s /snapshots/latest.snapshot /snapshots/\$current_epoch.epoch
    cp /snapshots/uploads/lachesis.snapshot /snapshots/tmp.snapshot
    cat /snapshots/latest.snapshot >> /snapshots/tmp.snapshot
    mv /snapshots/tmp.snapshot /snapshots/uploads/lachesis.snapshot
done
