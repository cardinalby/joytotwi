#!/bin/bash
echo "$HEROKU_API_KEY" | docker login --username=_ --password-stdin registry.heroku.com &&
docker push registry.heroku.com/$HEROKU_APP/web &&
imageId=$(docker inspect registry.heroku.com/$HEROKU_APP/web --format={{.Id}})
payload='{"updates":[{"type":"web","docker_image":"'"$imageId"'"}]}'
curl -n -X PATCH https://api.heroku.com/apps/$HEROKU_APP/formation \
-d "$payload" \
-H "Content-Type: application/json" \
-H "Accept: application/vnd.heroku+json; version=3.docker-releases" \
-H "Authorization: Bearer $HEROKU_API_KEY"