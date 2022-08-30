package idempotency

func AtMostOnce(handler interface{}, tableName string) interface{} {
	return WrapHandlerWithPersistence(handler)
}

// type Entry struct {
// 	Key         string `dynamodbav:"pk"`
// 	CurrentTime int64  `dynamodbav:"now"`
// 	ExpireAt    int64  `dynamodbav:"ttl"`
// }

// type Manager interface {
// 	SaveToDB(context.Context, Entry) error
// 	ReadFromDB(context.Context, string) (Entry, error)
// }

// type Service struct {
// 	table  *string
// 	client *dynamodb.Client
// }

// func (s *Service) SaveToDB(ctx context.Context, entry Entry) error {
// 	data, err := attributevalue.MarshalMap(entry)
// 	if err != nil {
// 		return errors.Wrap(err, "unable to write to database")
// 	}

// 	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
// 		TableName: s.table,
// 		Item:      data,
// 	})
// 	if err != nil {
// 		return errors.Wrap(err, fmt.Sprintf("unable to save entry to %s", *s.table))
// 	}

// 	fmt.Printf("%v saved to %s\n", data, *s.table)

// 	return nil
// }

// func (s *Service) ReadFromDB(ctx context.Context, pk string) (entry Entry, err error) {
// 	var data *dynamodb.GetItemOutput

// 	data, err = s.client.GetItem(ctx, &dynamodb.GetItemInput{
// 		TableName: s.table,
// 		Key: map[string]types.AttributeValue{
// 			"pk": &types.AttributeValueMemberS{Value: pk},
// 		},
// 	})
// 	if err != nil {
// 		return
// 	}

// 	fmt.Printf("items: %v\n", len(data.Item))

// 	err = attributevalue.UnmarshalMap(data.Item, &entry)
// 	if err != nil {
// 		return
// 	}

// 	if entry.ExpireAt <= time.Now().Unix() {
// 		fmt.Printf("EXPIRED entry: %#v\n", entry)
// 		entry = Entry{}
// 	} else {
// 		fmt.Printf("entry: %#v\n", entry)
// 	}

// 	return
// }

// func New(table string) Manager {
// 	cfg, _ := config.LoadDefaultConfig(context.Background())
// 	return &Service{
// 		table:  &table,
// 		client: dynamodb.NewFromConfig(cfg),
// 	}
// }
