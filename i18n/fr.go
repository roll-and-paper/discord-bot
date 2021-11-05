package i18n

var fr = D{
	"commands": D{
		"roll": D{
			"generic": D{
				"description": "Effectue un lancé de dé sur base d'une formule",
				"options": D{
					"expression": D{
						"description": "formule à executer",
					},
				},
			},
			"vampire-dark-ages": D{
				"description": "Effectue un lancé de dé sur le système Vampire Dark Ages",
				"options": D{
					"dices": D{
						"description": "Nombre de dés à lancer",
					},
					"difficulty": D{
						"description": "Difficulté du lancé",
					},
					"specialisation": D{
						"description": "Nombre de spécialité",
					},
				},
			},
		},
	},
	"messages": D{
		"welcome": "Hello.",
		"config":  "Voilà la configuration de ton serveur : ",
		"init": D{
			"start":                "Je commence à préparer ton serveur pour recevoir tes joueurs.",
			"players-role-created": "Le rôle 'Joueurs' a été créé",
			"master-role-created":  "Le rôle 'MJ' a été créé",
			"config-hint":          "Ce channel sera utilisé pour administrer tes parties et ton serveur (changement de config, ajout et suppression de PJ, demande de jet de dés, ...)",
			"finished":             "Le serveur est configuré, il ne te reste qu'à assigner ton MJ via la commande `{{.Config.Prefix}}set master @tonMj`",
		},
		"set": D{
			"master": D{
				"help": "Commande : {{.Config.Prefix}}set master `@master`",
			},
			"player": D{
				"help": "Command : {{.Config.Prefix}}set player `@joueur` nom du joueur",
			},
			"game-system": D{
				"help": "Commande : {{.Config.Prefix}}set game-system `value`\nValue : (vampire-dark-ages)",
			},
			"prefix": D{
				"help": "Command : {{.Config.Prefix}}set prefix `value`",
			},
		},
		"help": `Voila les differents commandes dispo : 
- **{{.Config.Prefix}}init** : Initialise le serveur, crée les rôles 'Joueurs' et 'MJ', crée les channels, ...
- **{{.Config.Prefix}}set master @masterName** : Attribue un membre au rôle de maitre de jeu
- **{{.Config.Prefix}}set player @playerName nom du personnage** : Attribue le rôle joueur, créer le channel associé et renomme le membre avec le nom du personnage
- **{{.Config.Prefix}}set game-system** : change le système de jeu du serveur (dispo : vampire-dark-ages)
- **{{.Config.Prefix}}set prefix newPrefix** : change le prefix utilisé par la valeur du nouveau prefix
- Toute autre commande est interprété comme une expression de dé, en fonction du système de jeu du serveur actuel.
`,
		"oops": "Ooops, une erreur c'est produite",
	},
	"name": D{
		"role": D{
			"players": "Joueurs",
			"master":  "MJ",
		},
		"channel": D{
			"players":  "Entres joueurs & maitres",
			"with-bot": "roll-and-paper-bot",
		},
	},
	"errors": D{
		"not-a-person":   "{{.Name}} ne correspond pas à une personne de ton serveur",
		"unknown-person": "je ne connais pas la personne {{.Name}}",
		"commands": D{
			"set": D{
				"bad-command":  "je ne vois pas ce que tu veux éditer, je ne connais pas d'éditeur pour '{{.Cmd}}', je ne suis capable de changer que {{.AllCmd}}",
				"unauthorized": "tu n'est pas authorisé à changer ma config",
				"master": D{
					"missing":      "tu n'as pas précisé ton maitre de jeu",
					"role-missing": "le role `MJ` n'exite pas, as tu initialisé le serveur avec `{{.Config.Prefix}}init` ?",
				},
				"player": D{
					"missing":                "tu n'as pas précisé ton joueur ou le nom du personnage",
					"role-missing":           "le role `Joueurs` n'exite pas, as tu initialisé le serveur avec `{{.Config.Prefix}}init` ?",
					"channel-missing":        "le channel `Entres joueurs & maitres` n'exite pas, as tu initialisé le serveur avec `{{.Config.Prefix}}init` ?",
					"have-already-character": "le joueur {{.Name}} interprete déjà le personnage {{.CharacterName}}, pour pouvoir lui réafecter un nouveau personnage, supprime le channel de son personne existant",
					"cannot-change-name":     "Je n'ai pas les droits pour changer le nom de {{.Name}} pour {{.CharacterName}}",
				},
				"game-system": D{
					"missing":             "tu n'as pas précisé le systeme de jeu désiré",
					"unknown-game-system": "je ne connais pas le systeme de jeu `{{.Wanted}}`, je ne connais que `{{.AllSystems}}`",
				},
				"prefix": D{
					"missing": "tu n'as pas précisé le nouveau préfix",
				},
			},
			"init": D{
				"unauthorized": "seul le propriétaire du serveur peu m'initialiser",
			},
			"roll": D{
				"unknown-game-system": "je n'ai pas de règles pour effectuer un jet sur le système choisi",
			},
			"unknown": "Commande inconnue : {0}",
		},
	},
}
