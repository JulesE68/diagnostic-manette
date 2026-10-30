# Diagnostic Manette

Logiciel Windows **portable** (un seul `.exe`, sans installation) qui teste l'usure d'une manette Xbox et donne une note sur 100, avec une interface façon Windows 7.

**Télécharger :** voir la page [Releases](../../releases/latest) et prendre `DiagnosticManette.exe`.

## Ce qu'il mesure
| Critère | Points | Mesure |
|---|---|---|
| Stick gauche / droit | 25 + 25 | Drift au repos, course maximale, régularité du contour |
| Gâchettes | 20 | Course, retour à zéro, progressivité |
| Boutons | 20 | Les 14 boutons répondent, doubles appuis parasites |
| Vibrations | 10 | Chaque moteur, confirmé par l'utilisateur |

Un **défaut critique** (stick bloqué hors du centre, drift important, zone morte, gâchette qui ne revient pas, bouton mort, moteur HS) est signalé en rouge et plafonne la note à 50. Le drift est aussi détecté en direct dès qu'un stick reste décentré au repos.

Le rapport s'enregistre en `.txt` à côté du programme, prêt à joindre à une annonce d'occasion.

## Utilisation
1. Lancer `DiagnosticManette.exe`. Au premier lancement, Windows SmartScreen peut avertir (programme non signé) : *Informations complémentaires* → *Exécuter quand même*.
2. Brancher la manette (USB ou Bluetooth) et appuyer sur un bouton.
3. Suivre les 7 étapes de l'assistant, puis ouvrir l'onglet *Rapport*.

Configuration requise : Windows 10 ou 11 avec Microsoft Edge WebView2 Runtime (déjà présent sur Windows 11). Aucune donnée n'est envoyée sur Internet.

**Limites :** le ressenti mécanique (bouton mou, stick qui accroche), la latence et l'état de la coque ne sont pas mesurés ; l'étape *Contrôle visuel* sert à les noter à part. Le bouton Xbox central n'est pas lisible.

## Structure du code
- `ui/src.html` : toute l'interface et la logique de test (HTML/CSS/JS). C'est le fichier à modifier.
- `build_ui.py` : intègre 7.css et l'icône dans `ui/app.html` (fichier embarqué dans le `.exe`).
- `main.go` : fenêtre Windows (WebView2) qui sert `ui/app.html` en local sur 127.0.0.1.
- `cmd/mkicon` : génère `rsrc_windows_amd64.syso` (icône et infos du `.exe`).

## Compiler
Prérequis : Go 1.24+, Python 3 avec Pillow, et une copie de [7.css](https://github.com/khang-nd/7.css) compilée (`7.css` de la branche `gh-pages`) placée dans `../deps/7css-pages/`.

```
python build_ui.py
go mod tidy
go run ./cmd/mkicon
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -H windowsgui" -o DiagnosticManette.exe .
```

`ui/app.html` et `rsrc_windows_amd64.syso` sont versionnés : `go build` suffit si l'interface n'a pas changé.

## Crédits
- Style : [7.css](https://github.com/khang-nd/7.css) de khang-nd (licence MIT)
- Fenêtre WebView2 : [go-webview2](https://github.com/jchv/go-webview2)
- Ressources Windows : [winres](https://github.com/tc-hib/winres)
