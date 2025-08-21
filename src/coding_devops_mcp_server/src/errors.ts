/**
 * 错误处理模块
 * 定义了系统中使用的所有自定义错误类型
 * 包括基础错误、API错误、认证错误等
 */

/**
 * 基础错误类
 * 所有自定义错误的基类，继承自 Error
 * 修复了 TypeScript 中 Error 类继承的原型链问题
 */
export class BaseError extends Error {
  constructor(message: string) {
    super(message);
    this.name = this.constructor.name;
    // Restore prototype chain in Node.js
    Object.setPrototypeOf(this, BaseError.prototype);
  }
}

/**
 * 配置错误类
 * 用于处理所有与配置相关的错误
 * 例如：配置文件格式错误、必需配置项缺失等
 */
export class ConfigurationError extends BaseError {
  constructor(message: string) {
    super(message);
    Object.setPrototypeOf(this, ConfigurationError.prototype);
  }
}

/**
 * API错误基类
 * 用于处理所有与API调用相关的错误
 * 包含HTTP状态码和原始响应数据
 * 作为所有API相关错误的基础类
 * @param message 错误信息
 * @param statusCode HTTP状态码，用于指示错误的具体类型
 * @param response 原始响应数据，保存服务器返回的完整信息
 */
export class ApiError extends BaseError {
  constructor(
    message: string,
    public readonly statusCode?: number,
    public readonly response?: unknown
  ) {
    super(message);
    Object.setPrototypeOf(this, ApiError.prototype);
  }
}

/**
 * 认证错误类
 * 用于处理所有与认证相关的错误
 * 例如：token无效、认证失败等
 * 默认返回401状态码
 */
/**
 * 认证错误类
 * 用于处理所有与认证相关的错误
 * 默认返回401状态码表示认证失败
 */
export class AuthenticationError extends ApiError {
  constructor(message: string = 'Authentication failed') {
    super(message, 401);
    Object.setPrototypeOf(this, AuthenticationError.prototype);
  }
}

/**
 * 资源未找到错误类
 * 用于处理请求的资源不存在的情况
 * 默认返回404状态码
 */
/**
 * 资源未找到错误类
 * 用于处理请求的资源不存在的情况
 * 默认返回404状态码表示资源未找到
 */
export class NotFoundError extends ApiError {
  constructor(message: string) {
    super(message, 404);
    Object.setPrototypeOf(this, NotFoundError.prototype);
  }
}

/**
 * MCP错误码枚举
 * 定义了系统中使用的标准错误码
 * - InvalidParams: 参数无效
 * - MethodNotFound: 方法未找到
 * - InternalError: 内部错误
 */
export enum ErrorCode {
  InvalidParams = 'invalid_params',
  MethodNotFound = 'method_not_found',
  InternalError = 'internal_error',
  ProjectCreateFailed = 'project_create_failed',
}

/**
 * MCP错误类
 * 系统核心错误类，用于处理标准化的错误情况
 * 包含错误码和错误信息
 */
export class McpError extends Error {
  constructor(public code: ErrorCode, message: string) {
    super(message);
    this.name = 'McpError';
    Object.setPrototypeOf(this, McpError.prototype);
  }
}