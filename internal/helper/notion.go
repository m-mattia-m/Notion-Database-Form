package helper

//type StoreObject struct {
//	Expiration     time.Time   `json:"expiration"`
//	RelevanceScore int         `json:"relevance_score"`
//	Object         interface{} `json:"object"`
//}
//type StoreDatabaseObject struct {
//	Expiration     time.Time       `json:"expiration"`
//	RelevanceScore int             `json:"relevance_score"`
//	Object         notion.Database `json:"object"`
//}
//type StorePageObject struct {
//	Expiration     time.Time   `json:"expiration"`
//	RelevanceScore int         `json:"relevance_score"`
//	Object         notion.Page `json:"object"`
//}
//
//func ConvertNotionCacheStringToObject(notionString string, notionObject interface{}) error {
//	var returnObject StoreObject
//	err := json.Unmarshal([]byte(notionString), &returnObject)
//	if err != nil {
//		return err
//	}
//
//	returnMap, ok := returnObject.Object.(map[string]interface{})
//	if !ok {
//		return fmt.Errorf("failed to convert cache-object to map")
//	}
//
//	objectType, ok := returnMap["object"]
//	if !ok {
//		return fmt.Errorf("failed to get 'object'-type from map")
//	}
//
//	if objectType == "database" {
//		database, err := convertNotionDatabaseStringToObject(notionString)
//		if err != nil {
//			return err
//		}
//		notionObject = database
//		return nil
//	}
//
//	if objectType == "page" {
//		page, err := convertNotionPageStringToObject(notionString)
//		if err != nil {
//			return err
//		}
//		notionObject = page
//		return nil
//	}
//	return nil
//}
//
//func convertNotionDatabaseStringToObject(notionString string) (*StoreDatabaseObject, error) {
//	var returnDatabaseObject *StoreDatabaseObject
//	err := json.Unmarshal([]byte(notionString), &returnDatabaseObject)
//	if err != nil {
//		return nil, err
//	}
//	return returnDatabaseObject, nil
//}
//
//func convertNotionPageStringToObject(notionString string) (*StorePageObject, error) {
//	var returnPageObject *StorePageObject
//	err := json.Unmarshal([]byte(notionString), &returnPageObject)
//	if err != nil {
//		return nil, err
//	}
//	return returnPageObject, nil
//}
