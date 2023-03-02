
:start
 :: Post Add Score
 	curl -Lvso /dev/null -d  "@GatewayLogin.json" -X POST http://localhost:8569/gateway/login
 	@REM curl GET http://localhost:8585/signup/49385234
 goto start

@REM curl -Lvso /dev/null -d  "@SignUp.json" -X POST http://localhost:8585/signUp
@Rem curl -Lvso /dev/null -d  "@GatewayLogin.json" -X POST http://localhost:8569/gateway/login