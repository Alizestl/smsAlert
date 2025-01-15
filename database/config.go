// Package database
// 使用viper第三方库，读取db的配置文件setting.yml的子节点，并对读取配置文件过程中发生的error进行处理
package database

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// 数据库配置项，全局变量，同一个包内的文件均可访问
var cfgDatabase *viper.Viper

func GenerateDatabaseConfig(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	if err := viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content)))); err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	// 数据库配置
	cfgDatabase = viper.Sub("settings.database") // 读取yaml 子节点settings.database
	if cfgDatabase == nil {
		panic("No found settings.database in the configuration")
	}
}

//// 格式化为 JSON
//func printAsJSON(cfg *viper.Viper) {
//	settings := cfg.AllSettings()
//
//	// 将配置序列化为 JSON
//	jsonData, err := json.MarshalIndent(settings, "", "  ")
//	if err != nil {
//		log.Fatalf("Failed to marshal cfgDatabase to JSON: %v", err)
//	}
//
//	fmt.Println("Configuration (JSON):")
//	fmt.Println(string(jsonData))
//}
//
//// 格式化为 YAML
//func printAsYAML(cfg *viper.Viper) {
//	settings := cfg.AllSettings()
//
//	// 将配置序列化为 YAML
//	yamlData, err := yaml.Marshal(settings)
//	if err != nil {
//		log.Fatalf("Failed to marshal cfgDatabase to YAML: %v", err)
//	}
//
//	fmt.Println("Configuration (YAML):")
//	fmt.Println(string(yamlData))
//}
