package relojModel

import (
	"../../util/config"
	"../relojEntitie"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vjeantet/jodaTime"
	"math"
	"strconv"
	"time"
)

var Db *sql.DB

func InitModel(configFile string) {
	var err error

	Db, err = sql.Open("mysql",
		config.GenerateStringMysqlDbConnection(configFile))
	if err != nil {
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

//func FindById(id int) (user relojEntitie.User, err error) {
//	user = relojEntitie.User{}
//	err = Db.QueryRow("select id, username, password from users where id = $1",
//		id).Scan(&user.Id, &user.Username, &user.Password)
//	return user, nil
//}
//
//func FindByUsernameAndPassword(username string, password string) (user relojEntitie.User, err error) {
//	user = relojEntitie.User{}
//	err = Db.QueryRow("select id, username, password from users where username = $1 and password = $2",
//		username, password).Scan(&user.Id, &user.Username, &user.Password)
//	return user, nil
//}

func FindAll() (relojs []relojEntitie.Reloj, err error) {
	rows, err := Db.Query("select cod_us, ing, sale, estado, tiempo from reloj limit 20")
	for rows.Next() {
		reloj := relojEntitie.Reloj{}
		err = rows.Scan(&reloj.CodUs, &reloj.Ing, &reloj.Sale, &reloj.Estado, &reloj.Tiempo)
		relojs = append(relojs, reloj)
	}
	return relojs, nil
}

func VerifyDni(dni string) (hasError bool, err error) {
	user := relojEntitie.User{}
	err = Db.QueryRow("SELECT nro_doc, psw FROM usuario WHERE nro_doc = ?", dni).Scan(&user.NroDoc, &user.Psw)
	if err != nil {
		return true, err
	}

	return false, err
}

func VerifyDniAndPassword(user relojEntitie.User) (hasError bool, err error) {
	h := md5.New()
	h.Write([]byte(user.Psw))
	bs := h.Sum(nil)
	pswToString := hex.EncodeToString(bs)

	err = Db.QueryRow("SELECT nro_doc, psw FROM usuario WHERE nro_doc = ? AND psw = ?", user.NroDoc, pswToString).Scan(&user.NroDoc, &user.Psw)
	if err != nil {
		return true, err
	}

	return false, err
}

func VerifyMarkByOperation(nroDoc string, operation string) (hasError bool, err error, message string) {
	hasError = false
	var notificacion = false
	var diferencia float64

	user := relojEntitie.User{}
	err = Db.QueryRow("SELECT * FROM reloj_horarios as hor,usuario as us WHERE hor.codigo=? and hor.cod_us=us.cod_us", nroDoc).Scan(&user.NroDoc, &user.CodUs, &user.Nombres, &user.Apellido, &user.In1, &user.In2, &user.Flexible)
	if err != nil {
		return true, err, "Error en validacion de tarjeta!"
	} else {
		var relojCuentaRow = Db.QueryRow("SELECT COUNT(*) FROM reloj WHERE cod_us=? and estado='I' and revision='N'", user.CodUs)
		var relojCuentaCount = 0
		_ = relojCuentaRow.Scan(&relojCuentaCount)

		var relojErrorRow = Db.QueryRow("SELECT COUNT(*) FROM reloj WHERE cod_us=? and revision='E'", user.CodUs)
		var relojErrorCount = 0
		_ = relojErrorRow.Scan(&relojErrorCount)

		var horaMarca = jodaTime.Format("H:m:s", time.Now())
		var usuarioReloj = user.CodUs
		var usuarioNombre = user.Nombres + " " + user.Apellido

		if relojErrorCount != 0 {
			message = usuarioNombre + ", usted tiene un error anterior que el encargado de personal notificara en breve. Puede seguir marcando ya que todas las marcaciones son tomadas."

			_, err := Db.Exec("INSERT INTO reloj (cod_us,revision,ing,sale,estado,tiempo) VALUES (?,'E',CURRENT_TIMESTAMP,'','I','0')", usuarioReloj)
			if err != nil {
				return true, err, "SQL ERROR"
			}

			hasError = true
		} else {

			if relojCuentaCount != 0 {
				if operation == "I" {
					message = usuarioNombre + ", no marco el/la ingreso/salida anterior (El encargado de personal lo notificara por mail del error observado). Todas las marcaciones que realice a posterior las tomara de igual manera"

					_, err := Db.Exec("UPDATE reloj SET revision='E' WHERE cod_us=? AND estado='I'", usuarioReloj)
					if err != nil {
						return true, err, "SQL ERROR"
					}

					_, err = Db.Exec("INSERT INTO reloj (cod_us,revision,ing,sale,estado,tiempo) VALUES (?,'E',CURRENT_TIMESTAMP,'0','I','0')", usuarioReloj)
					if err != nil {
						return true, err, "SQL ERROR"
					}

					hasError = true
				} //mal
				if operation == "O" {
					message = usuarioNombre + " marca salida a las: " + horaMarca

					_, err := Db.Exec("UPDATE reloj SET  sale=CURRENT_TIMESTAMP, estado='O',tiempo=SEC_TO_TIME((UNIX_TIMESTAMP())-(UNIX_TIMESTAMP(ing))) WHERE cod_us=? AND estado='I' AND revision='N'", usuarioReloj)
					if err != nil {
						return true, err, "SQL ERROR"
					}

					hasError = false
				} //ok
			}

			if relojCuentaCount == 0 {
				if operation == "I" { // este es un ingreso exitoso

					/////// VOY A BUSCAR SI TIENE UN INGRESO EN EL DIA
					var relojMarcacionesDiariasRow = Db.QueryRow("SELECT COUNT(*) FROM reloj WHERE cod_us=? AND date(ing)=date(now())", user.CodUs)
					var relojMarcacionesDiariasCount = 0
					_ = relojMarcacionesDiariasRow.Scan(&relojMarcacionesDiariasCount)

					if relojMarcacionesDiariasCount == 0 { // es el primero del dia
						diferencia = minutos_transcurridos(horaMarca, user.In1)

						if diferencia >= 10 {
							notificacion = true
						}
					} else { // es la segunda del dia
						diferencia = minutos_transcurridos(horaMarca, user.In2)

						if diferencia >= 0 {
							notificacion = true
						}
					}

					message = usuarioNombre + " marca ingreso a las: " + horaMarca

					_, err = Db.Exec("INSERT INTO reloj (cod_us,revision,ing,sale,estado,tiempo) VALUES (?,'N',CURRENT_TIMESTAMP,'0','I','0')", usuarioReloj)
					if err != nil {
						return true, err, "SQL ERROR"
					}

					hasError = false

				} //ok

				if operation == "O" {
					message = usuarioNombre + ", no marco el/la ingreso/salida anterior (El encargado de personal lo notificara por mail del error observado). Todas las marcaciones que realice a posterior las tomar&aacute; de igual manera"

					_, err = Db.Exec("INSERT INTO reloj (cod_us,revision,ing,sale,estado,tiempo) VALUES (?,'E','0',CURRENT_TIMESTAMP,'O','0')", usuarioReloj)
					if err != nil {
						return true, err, "SQL ERROR"
					}

					hasError = true
				} //mal
			}
		}

		enviarNotificacionSupervisor(usuarioReloj, horaMarca, usuarioNombre)

		if notificacion && user.Flexible == 0 {
			var area = relojEntitie.Area{}
			var diferenciaparanoti string

			_ = Db.QueryRow("SELECT a.supervisor,a.encargado FROM usu_x_area as us INNER JOIN areas AS a ON a.id_area=us.id_area WHERE us.cod_us=?", usuarioReloj).Scan(&area.Supervisor, &area.Encargado)

			if diferencia < 60 {
				diferenciaparanoti = strconv.FormatFloat(diferencia, 'f', 5, 64) + " minutos"
			} else {
				var horas = math.Round(diferencia / 60)
				var minutosDeHoras = horas * 60
				var minutosRestantes = diferencia - minutosDeHoras
				diferenciaparanoti = strconv.FormatFloat(horas, 'f', 5, 64) + " hora(s) " + strconv.FormatFloat(minutosRestantes, 'f', 5, 64) + " minuto(s)"

				var mensajeUsuario = utf8_decode("Se le notifica que el día de la fecha  ud posee una llegada tarde de " + diferenciaparanoti + ", marcando ingreso a las " + horaMarca + ", pudiendo visualizarlo en la funcionalidad Consulta de Reloj. ")
				var mensajeEncargado = utf8_decode("Se notifica que el colaborador " + usuarioNombre + " en día de la fecha posee llegada tarde de " + diferenciaparanoti + ", marcando ingreso a las " + horaMarca + ", pudiendo visualizarlo en la Consulta Reloj x Area.")

				var botonFuncionalidad = "<a href=\"administracion.php?usuario=_esteusuariosedebecambiar_&sesion=_estasesionsedebecambiar_&pagina=empl/reloj_lista_usuario.php\" class=\"btn btn-success btn-mensaje\" >Consulta Reloj</a>"
				_, _ = Db.Exec("INSERT INTO mensajes (id, titulo, introtext, texto, img, origen, nombre_origen, destino, tipo_mensaje_id, categoria_mensaje_id, funcionalidad_usuario, funcionalidad_id, funcionalidad_vencimiento, destacado, fecha_alta, fecha_baja, fecha_modificacion, fecha_lectura) VALUES (NULL, '', 'Llegada tarde de ?', '<p><strong>Llegada tarde: <br></strong>? </p><p>&nbsp;</p><p>Muchas gracias por su atenci&oacute;n,</p><p>Intranet Vanesa Dur&aacute;n</p> <p class=\"text-left\">?</p>', NULL, 'Llegada tarde', 'Intranet VD', '? mensaje', '1', '5', NULL, NULL, NULL, '0', NOW(), NULL , '', NULL )", usuarioReloj, mensajeUsuario, botonFuncionalidad, usuarioReloj)

				var usuarioMensaje string
				if area.Encargado == usuarioReloj {
					usuarioMensaje = area.Supervisor
				} else {
					usuarioMensaje = area.Encargado
				}

				botonFuncionalidad = "<a href=\"administracion.php?usuario=_esteusuariosedebecambiar_&sesion=_estasesionsedebecambiar_&pagina=empl/reloj_usuario_lista.php\"class=\"btn btn-success btn-mensaje\" >Consulta Reloj x Area</a>"
				_, _ = Db.Exec("INSERT INTO mensajes (id, titulo, introtext, texto, img, origen, nombre_origen, destino, tipo_mensaje_id, categoria_mensaje_id, funcionalidad_usuario, funcionalidad_id, funcionalidad_vencimiento, destacado, fecha_alta, fecha_baja, fecha_modificacion, fecha_lectura) VALUES (NULL, '', 'Llegada tarde de ?', '<p><strong>Llegada tarde: <br></strong>? </p><p>&nbsp;</p><p>Muchas gracias por su atenci&oacute;n,</p><p>Intranet Vanesa Dur&aacute;n</p> <p class=\"text-left\">?</p>', NULL, 'Llegada tarde', 'Intranet VD', '?', '1', '5', NULL, NULL, NULL, '0', NOW(), NULL , '', NULL )", usuarioReloj, mensajeEncargado, botonFuncionalidad, usuarioMensaje)

				botonFuncionalidad = "<a href=\"administracion.php?usuario=_esteusuariosedebecambiar_&sesion=_estasesionsedebecambiar_&pagina=empl/reloj_usuario_lista.php\"class=\"btn btn-success btn-mensaje\" >Consulta Reloj x Area</a>"
				_, _ = Db.Exec("INSERT INTO mensajes (id, titulo, introtext, texto, img, origen, nombre_origen, destino, tipo_mensaje_id, categoria_mensaje_id, funcionalidad_usuario, funcionalidad_id, funcionalidad_vencimiento, destacado, fecha_alta, fecha_baja, fecha_modificacion, fecha_lectura) VALUES (NULL, '', 'Llegada tarde de ?', '<p><strong>Llegada tarde: <br></strong>? </p><p>&nbsp;</p><p>Muchas gracias por su atenci&oacute;n,</p><p>Intranet Vanesa Dur&aacute;n</p> <p class=\"text-left\">?</p>', NULL, 'Llegada tarde', 'Intranet VD', '? mensaje', '1', '5', NULL, NULL, NULL, '0', NOW(), NULL , '', NULL )", usuarioReloj, mensajeEncargado, botonFuncionalidad, "RODRIGO")
			}
		}

		return hasError, err, message
	}
}

func enviarNotificacionSupervisor(usuarioReloj string, horaMarca string, usuarioNombre string) {

	var area = relojEntitie.Area{}
	_ = Db.QueryRow("SELECT a.supervisor,a.encargado FROM usu_x_area as us INNER JOIN areas AS a ON a.id_area=us.id_area WHERE us.cod_us=?", usuarioReloj).Scan(&area.Supervisor, &area.Encargado)

	var usuarioRRHH = "RODRIGO"

	var mensajeUsuario = utf8_decode("Se le notifica que el día de la fecha  ud marcó a las " + horaMarca + " desde fuera de la empresa, pudiendo visualizarlo en la funcionalidad Consulta de Reloj. ")
	var mensajeAutoridad = utf8_decode("Se notifica que el colaborador " + usuarioNombre + " en día de la fecha marcó a las " + horaMarca + " desde fuera de la empresa, pudiendo visualizarlo en la Consulta Reloj x Area. ")

	_, _ = Db.Exec("INSERT INTO mensajes (id, titulo, introtext, texto, img, origen, nombre_origen, destino, tipo_mensaje_id, categoria_mensaje_id, funcionalidad_usuario, funcionalidad_id, funcionalidad_vencimiento, destacado, fecha_alta, fecha_baja, fecha_modificacion, fecha_lectura)  VALUES (NULL, '', 'Marcacion externa de ?', '<p><strong>Marcaci&oacute;n externa: <br></strong>? </p> <p>&nbsp;</p> <p>Muchas gracias por su atenci&oacute;n,</p> <p>Intranet Vanesa Dur&aacute;n</p>', NULL, 'Marcacion externa', 'Intranet VD', '?', '1', '5', NULL, NULL, NULL, '0', NOW(), NULL , '', NULL )", usuarioReloj, mensajeUsuario, usuarioReloj)
	_, _ = Db.Exec("INSERT INTO mensajes (id, titulo, introtext, texto, img, origen, nombre_origen, destino, tipo_mensaje_id, categoria_mensaje_id, funcionalidad_usuario, funcionalidad_id, funcionalidad_vencimiento, destacado, fecha_alta, fecha_baja, fecha_modificacion, fecha_lectura)  VALUES (NULL, '', 'Marcacion externa de ?', '<p><strong>Marcaci&oacute;n externa: <br></strong>? </p> <p>&nbsp;</p> <p>Muchas gracias por su atenci&oacute;n,</p> <p>Intranet Vanesa Dur&aacute;n</p>', NULL, 'Marcacion externa', 'Intranet VD', '?', '1', '5', NULL, NULL, NULL, '0', NOW(), NULL , '', NULL )", usuarioReloj, mensajeAutoridad, area.Encargado)
	_, _ = Db.Exec("INSERT INTO mensajes (id, titulo, introtext, texto, img, origen, nombre_origen, destino, tipo_mensaje_id, categoria_mensaje_id, funcionalidad_usuario, funcionalidad_id, funcionalidad_vencimiento, destacado, fecha_alta, fecha_baja, fecha_modificacion, fecha_lectura)  VALUES (NULL, '', 'Marcacion externa de ?', '<p><strong>Marcaci&oacute;n externa: <br></strong>? </p> <p>&nbsp;</p> <p>Muchas gracias por su atenci&oacute;n,</p> <p>Intranet Vanesa Dur&aacute;n</p>', NULL, 'Marcacion externa', 'Intranet VD', '?', '1', '5', NULL, NULL, NULL, '0', NOW(), NULL , '', NULL )", usuarioReloj, mensajeAutoridad, usuarioRRHH)
}

//func Create(user relojEntitie.User) (id int, err error) {
//	statement := "insert into users (username, password) values ($1, $2) returning id"
//	stmt, err := Db.Prepare(statement)
//	if err != nil {
//		return
//	}
//	defer stmt.Close()
//	err = stmt.QueryRow(user.Username, user.Password).Scan(&user.Id)
//	return user.Id, err
//}
//
//func Update(user relojEntitie.User) (err error) {
//	_, err = Db.Exec("update users set username = $2, password = $3 where id = $1",
//		user.Id, user.Username, user.Password)
//	return
//}
//
//func Delete(id int) (err error) {
//	_, err = Db.Exec("delete from users where id = $1", id)
//	return
//}

func minutos_transcurridos(date1 string, date2 string) (diference float64) {
	var horaMarcaTime, _ = jodaTime.Parse("H:m:s", date1)
	var horaMarcaTimeIn, _ = jodaTime.Parse("H:m:s", date2)
	return horaMarcaTime.Sub(horaMarcaTimeIn).Minutes()
}

func utf8_decode(str string) string {
	var result string
	for i := range str {
		result += string(str[i])
	}
	return result
}
