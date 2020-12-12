package response

var mapping map[string]string

func Initialize() {
	mapping = make(map[string]string)

	mapping["1001"] = "The admin profile is initialized"

}
