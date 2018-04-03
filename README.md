# Upsilon Garden

Ceci est un projet d'essai d'Upsilon Garden en Go. 

Objectif: 

* Liste des Jardins
* Ajouter un nouveau jardin.
* Acceder a un jardin
* Supprimer un jardin.
* Voir le contenu du jardin ( parcelles & plantes en cours de production )
* Ajouter une graine au jardin. (parmis une selection fixe)
* Observer humidité du sol (dur <0.2 , sec < 0.4, normal < 0.6, humide < 0.8, boueux < 0.95, submergé < 1)
* Observer aspect de la plante ( savoir si elle est "en bonne santée", "faiblarde", "mourante")
* Arroser une parcelle. 
* Assurer que si une plante est trop arrosée (ou trop peu) elle dépérie
* Récolter une plante ( retire la plante du jardin )

Requetes: 

* CRUD: /gardens 					Gestion des jardins.
* CRUD: /gardens/:id/plants 			Gestion des plantes dans un jardin. 
* GET:  /gardens/:id/hydro/:parcel	Action: Récuperer l'hydro courante de la parcelle ciblée; MAJ Jardin 
	( GET /gardens/:id dois maintnenat retourné cette valeur d'hydro pour la parcelle.)
* PUT:	/gardens/:id/hydro/:parcel	Action: Arrose la parcelle 

Sur chaque GET /gardens et /gardens/:id/plants/ on dois mettre a jour le jardin (avancement des tours et evolutions des plantes)

Create Garden: POST /gardens : { name: }
Update Garden: PUT /gardens/:id : { name: }

Create Plant: POST /gardens/:id/plants : { name: , plant_type: , parcel: }
Update Plant: PUT /gardens/:id/plants/:id { name: }

Water Parcel: POST /gardens/:id/hydro/:parcel {power:}

Objets: 

Garden : 
{
	id,
	name,
	parcels: [Parcel],
	last_update,
	next_update,
	plants: [Plant]
}

Parcel : 
{
	id,
	position,
	hydro,
	plant_id,
	base_hydro,
	running_hydro: [HydroEvent]
	next_hydro_end
}

Plant 
{
	id,
	level,
	name,
	plant_type,
	target_hydro: HydroRange,
	next_update,
	sp_per_level,			// 5
	sp_max,					// commence a 50
	sp_current				// commence a 25
}

HydroEvent
{
	begin_date,
	end_date,
	power
}

HydroRange 
{
	min_not_dead,
	max_not_dead,
	min_ok,
	max_ok,
	min_super,
	max_super,
}

Mise a jour du jardin: 

La projection indique les modifications de la plante sur le temps a partir de Now jusqu'au prochain evenement;

Les evenements a prendre en compte sont: 
* HydroEvent ( Fin d'effet d'un arrosage ) => On retire l'event d'hydro.
* Plant#Update ( La plante change de niveau ) => On augmente le niveau de 1 et on ajoute sp_per_level a sp_current et sp_max

Les modifications sur la plante sont :

* Calcule de l'hydro de la parcelle: min( Parcel#base_hydro + Sum(Parcel#HydroEvent#power), 1 )
* Si l'hydro est dans le range super: + 3 sp / tour 
* Si l'hydro est dans le range ok: + 1 sp / tour
* si l'hydro est dans le range not_dead : -3 sp / tour 
* si l'hydro est en dehors du range not_dead: -20 sp / tour

Si les SP (structure points) de la plante tombent a 0; la plante est automatiquement supprimé du jardin. 

# Configuration

Afin de se connecter a une DB postgres, les variables 

* config.DB_USER
* config.DB_PASSWORD
* config.DB_NAME
* config.DB_HOST

doivent etre créer. Un fichier config.go.sample est pret a l'usage, il suffit de le copier vers config.go et d'y remplacer les valeurs par des valeurs approprié. 

