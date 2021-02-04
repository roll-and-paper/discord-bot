package i18n

var fr = D{
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
		},
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
			},
			"init": D{
				"unauthorized": "seul le propriétaire du serveur peu m'initialiser",
			},
			"unknown": "Commande inconnue : {0}",
		},
	},
}
