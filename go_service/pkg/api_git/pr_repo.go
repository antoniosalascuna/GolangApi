package api_git

import (
	"context"
	"fmt"
	"go_service/tools"
	"time"

	"github.com/go-git/go-git/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
*Funcion que toma como parametros los datos que vienen
del la ruta para tomar y retornar los cambios que se han realizado en el
Pull request
*/
/* func pullRequest(pullRequest PRCreate) (*PRCreate, error) {
	var diffRequest = DiffRequest{
		pullRequest.UrlRepoReceivePR,
		"",
		pullRequest.UrlRepoCreatePR,
		pullRequest.CommitHash,
	}
	//Funcion que comprueba las diferencias que hay la rama main
	// con la rama que esta actualizando
	changes, err := diffTreeRepos(diffRequest)
	if err == nil {
		//Guarda en el atributo Path, la accion que se esta
		//realizando dentro del repositorio principal
		pullRequest.Patch = changes.String()
		return &pullRequest, err
	}
	if err != nil {
		return nil, err
	}

	return &pullRequest, err
} */

func pullRequestV2(pullRequest PRCreate) (*PRCreate, error) {
	var diffRequest = DiffRequestBranches{
		pullRequest.RepoName,
		pullRequest.BaseBranchNamePR,
		pullRequest.RefBranchNamePR,
		pullRequest.ProjecName,
	}

	//Funcion que comprueba las diferencias que hay la rama main
	// con la rama que esta actualizando
	changes, err := diffTreeBranches(diffRequest)
	lenchange := changes.Len()

	if err == nil {
		//Guarda en el atributo Path, la accion que se esta
		//realizando dentro del repositorio principal

		for _, value := range changes {

			patch, err := value.Patch()

			//from, to, err := c.Files()

			if err != nil {
				return nil, err
			}

			mekrlitre, err := value.Action()

			if err != nil {
				return nil, err
			}

			pullRequest.ActionsAndPaths = append(pullRequest.ActionsAndPaths, ActionPaths{
				Action:       string(mekrlitre.String()),
				Path:         patch.Stats().String(),
				PatchContent: patch.String(),
			})
		}

		GetCountPr, err := GetPRByRepoNameAndProjectName(pullRequest.RepoName, pullRequest.ProjecName)

		if err != nil {
			return nil, err
		}

		pullRequest.NumberPR = len(GetCountPr) + 1
		pullRequest.Patch = changes.String()
		pullRequest.LenChanges = lenchange

		return &pullRequest, err
	}
	if err != nil {
		return nil, err
	}

	return &pullRequest, err
}

func GetPRByRepoNameAndProjectName(RepoName string, ProjecName string) ([]*PRCreate, error) {

	var Pr []*PRCreate

	if RepoName == "" || ProjecName == "" {
		return Pr, git.ErrRepositoryNotExists
	} else {
		DBclient := tools.ConnectionDB()

		database := DBclient.Database("go_git")

		filter := bson.D{
			{"$and",
				bson.A{
					bson.D{{
						Key:   "RepoName",
						Value: RepoName,
					}},
					bson.D{{
						Key:   "ProjecName",
						Value: ProjecName,
					}},
				},
			},
		}

		PRcollection := database.Collection("PR_collection")

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		result, err := PRcollection.Find(ctx, filter)

		for result.Next(context.TODO()) {
			var elem PRCreate
			err := result.Decode(&elem)

			if err != nil {
				return Pr, err
			}
			Pr = append(Pr, &elem)
		}

		if err := result.Err(); err != nil {
			return Pr, err
		}

		result.Close(context.TODO())

		fmt.Println(len(Pr))

		return Pr, err
	}

}

/*
* Function that get a record from the collection
* through the {id} as parameter
 */
func GetOnePr(id primitive.ObjectID) (*PRCreate, error) {

	var Pr PRCreate

	//Open Connection with MongoDD
	//Funtion locate in tools Directory
	DBclient := tools.ConnectionDB()

	filter := bson.D{{Key: "_id", Value: id}}

	database := DBclient.Database("go_git")

	PRcollection := database.Collection("PR_collection")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// FindOne Query from MongoDB
	//This Query receives 2 parameters
	// The Firts is the context
	// The Second is the Filter (in this case is the {id} variable)

	err := PRcollection.FindOne(ctx, filter).Decode(&Pr)

	//RepoName := Pr.RepoName
	//ProjectName := Pr.ProjecName
	//RefBranchName := Pr.RefBranchNamePR
	//	BaseBrancName := Pr.BaseBranchNamePR
	//	var CommitArrayResult []CommitInfo

	//fmt.Println(CommitArrayResult)

	//Pr.CommitsHistory = CommitHistory

	return &Pr, err

	///Disconnect from Database
	//DBclient = tools.Disconnect()

}

func InsertOnePr(newPR PRCreate) (*mongo.InsertOneResult, error) {

	DBclient := tools.ConnectionDB()

	database := DBclient.Database("go_git")

	PRcollection := database.Collection("PR_collection")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := PRcollection.InsertOne(ctx, newPR)

	if err != nil {
		return result, err
	}

	insertedid := result.InsertedID

	fmt.Println(insertedid)

	//PRcollection.FindOne(ctx, PRCreate{_id: insertedid}).Decode(&Pr)

	tools.Disconnect(DBclient, ctx)

	return result, err

}

func UpdateStatusPR(id primitive.ObjectID) (*mongo.UpdateResult, error) {

	filter := bson.D{{"_id", id}}

	update := bson.D{{"$set", bson.D{{"IsOpen", false}}}}

	DBclient := tools.ConnectionDB()

	database := DBclient.Database("go_git")

	PRcollection := database.Collection("PR_collection")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := PRcollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return result, err
	}

	tools.Disconnect(DBclient, ctx)

	return result, err

}
