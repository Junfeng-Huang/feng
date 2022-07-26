# feng
A mini Web framework for Go. Its name is feng.

Compile build_feng_tool/feng_tool to generate the tool for feng.

feng contains five modules:
1. It is used to encapsulate common methods of request and response, transfer data between goroutines, and control the Context of goroutines.
2. It is used for HTTP method matching, static routing matching, dynamic routing matching, and routing of batch common prefix settings.
3. The middleware mechanism that can call middleware in the framework.
4. Services and service containers that can manage the relationship between modules and reduce the coupling between modules.
5. Command line tool with application management commands and debug mode



