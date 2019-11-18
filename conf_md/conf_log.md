
**配置说明**
{
    "TimeFormat":"2006-01-02 15:04:05", // 输出日志开头时间格式(详见**时间格式**)
    "Console": {            // 控制台日志配置
        "level": "TRAC",    // 控制台日志输出等级(详见**日志等级**)
        "color": true       // 控制台日志颜色开关 
    },
    "File": {                   // 文件日志配置
        "filename": "out.log",  // 初始日志文件名
        "level": "TRAC",        // 日志文件日志输出等级
        "daily": true,          // 跨天后是否创建新日志文件，当append=true时有效
        "maxlines": 1000000,    // 日志文件最大行数，当append=true时有效
        "maxsize": 1,           // 日志文件最大大小，当append=true时有效
        "maxdays": 3,           // 日志文件有效期
        "append": true,         // 是否支持日志追加
        "permit": "0660"        // 新创建的日志文件权限属性
    },
}

**日志等级**
等级	配置	释义	                                        控制台颜色
0	    EMER	系统级紧急，比如磁盘出错，内存异常，网络不可用等	红色底
1	    ALRT	系统级警告，比如数据库访问异常，配置文件出错等	    紫色
2	    CRIT	系统级危险，比如权限出错，访问异常等	            蓝色
3	    EROR	用户级错误	                                    红色
4	    WARN	用户级警告	                                    黄色
5	    INFO	用户级重要	                                    天蓝色
6	    DEBG	用户级调试	                                    绿色
7	    TRAC	用户级基本输出	                                绿色

**时间格式**
ANSIC	        "Mon Jan _2 15:04:05 2006"
UnixDate	    "Mon Jan _2 15:04:05 MST 2006"
RubyDate	    "Mon Jan 02 15:04:05 -0700 2006"
RFC822	        "02 Jan 06 15:04 MST"
RFC822Z	        "02 Jan 06 15:04 -0700"
RFC850	        "Monday, 02-Jan-06 15:04:05 MST"
RFC1123	        "Mon, 02 Jan 2006 15:04:05 MST"
RFC1123Z	    "Mon, 02 Jan 2006 15:04:05 -0700"
RFC3339	        "2006-01-02T15:04:05Z07:00"
RFC3339Nano	    "2006-01-02T15:04:05.999999999Z07:00"
Kitchen	        "3:04PM"
Stamp	        "Jan _2 15:04:05"
StampMilli	    "Jan _2 15:04:05.000"
StampMicro	    "Jan _2 15:04:05.000000"
StampNano	    "Jan _2 15:04:05.000000000"
RFC3339Nano1    "2006-01-02 15:04:05.999999999 -0700 MST"
DEFAULT	        "2006-01-02 15:04:05"