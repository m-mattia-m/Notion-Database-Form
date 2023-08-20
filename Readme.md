# Notion Forms

## Environment-variables

Go to [environment-variables](/.env]).



## Special

| Term | Description  |
|--|---|
| Time format | RFC3339  |
|  |   |
|  |   |



## ToDo

- Authentifiziere ob der User Zugriff auf die Datenbank/Form hat indem die im form/db-Objekt die Notion-Owner-Id eingetragen wird, danach kann mit dem user-objekt die owner-id abgeglichen werden (ggf. auch anstatt notion-owner-id die iam-id -> diese könnte ohne zwischenschritt überrüft werden)
- https://github.com/sirupsen/logrus einbauen im logger für Farben und timestaps -> ggf. kann auch kategorisiert werden -> warning, fatal, error, ...