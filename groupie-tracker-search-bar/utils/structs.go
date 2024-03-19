package test

type LocationIndex struct {
	Index []LocationData `json:"index"`
}

type DateIndex struct {
	Index []DatesData `json:"index"`
}

type RelationIndex struct {
	Index []DatesLocationsData `json:"index"`
}

type DatesLocationsData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type DatesData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type LocationData struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Artist2 struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Location LocationData
	Date DatesData
	Relation DatesLocationsData
}
