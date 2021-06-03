/**
 * Created by Wangwei on 2019-05-31 20:46.
 */

package conf

type AppConfig struct {
	JwtSecret    string `dsn:"jwt_secret"`
	PwdMd5Secret string `dsn:"pwd_md5_secret"`
	//JwtSecret         string `toml:"jwt_secret"`
	//PwdMd5Secret      string `toml:"pwd_md5_secret"`
	SuperAdminPassord string `toml:"super_admin_password"`
	SuperAdminAccount string `toml:"super_admin_account"`
	SuperAdminPhone   string `toml:"super_admin_phone"`
	DefaultPassword   string `toml:"default_password"`
}
