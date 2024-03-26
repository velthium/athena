#!/bin/bash

# Exécution des scripts SQL avec les nouveaux paramètres de connexion

psql -U athena_user -d athena_db -a -f 00-cosmos.sql
psql -U athena_user -d athena_db -a -f 00-fees.sql
psql -U athena_user -d athena_db -a -f 01-profiles.sql
psql -U athena_user -d athena_db -a -f 02-subspaces.sql
psql -U athena_user -d athena_db -a -f 03-relationships.sql
psql -U athena_user -d athena_db -a -f 04-posts.sql
psql -U athena_user -d athena_db -a -f 05-reports.sql
psql -U athena_user -d athena_db -a -f 06-reactions.sql
psql -U athena_user -d athena_db -a -f 07-contracts.sql
psql -U athena_user -d athena_db -a -f 08-tips.sql
psql -U athena_user -d athena_db -a -f 09-authz.sql
psql -U athena_user -d athena_db -a -f 10-feegrant.sql
psql -U athena_user -d athena_db -a -f 11-notifications.sql
psql -U athena_user -d athena_db -a -f 12-profile-score.sql
psql -U athena_user -d athena_db -a -f 13-profile-functions.sql
psql -U athena_user -d athena_db -a -f 14-profile-counters.sql