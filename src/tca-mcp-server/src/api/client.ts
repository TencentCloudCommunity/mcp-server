import { createHash } from 'crypto';
import * as ini from 'ini';
import * as fs from 'fs';

export interface CommonParams {
    orgSid: string;
    teamName: string;
    repoId?: string;
    projectId?: string;
}

export interface TcaClientOptions {
    token: string;
    userName: string;
    baseUrl: string;
}

export interface TcaClientHeaders {
    "TCA-USERID": string;
    "TCA-TICKET": string;
    "TCA-TIMESTAMP": string;
}

type RequestOptions = {
    method: string;
    headers: Record<string, string>;
    body?: any;
}

export default class TcaClient {
    private static instance: TcaClient | null = null;
    private static commonParams: CommonParams | null = null;

    private token: string;
    private userName: string;
    public baseUrl: string;

    constructor(options: TcaClientOptions) {
        this.token = options.token;
        this.userName = options.userName;
        this.baseUrl = options.baseUrl;
    }
    static initTcaClient(options: TcaClientOptions) {
        if (!TcaClient.instance) {
            TcaClient.instance = new TcaClient(options);   
        }
    }
    static getInstance(): TcaClient {
        if (TcaClient.instance) {
            return TcaClient.instance;
        } else {
            throw new Error("TcaClient is not initialized");
        }
    }

    static async initConfig(configPath: string) {
        // 解析 ini 配置文件
        const data = await fs.promises.readFile(configPath, 'utf-8');
        const config = ini.parse(data);
        const configSection = config.config || {};
        TcaClient.commonParams = {
            orgSid: configSection.org_sid,
            teamName: configSection.team_name,
            repoId: configSection.repo_id,
            projectId: configSection.project_id,
        }
        return 
    }
    static getcommonParams(): CommonParams{
        if (!TcaClient.commonParams) {
            throw new Error("common params is not initialized");
        }
        return TcaClient.commonParams;
    }

    getHeader(): TcaClientHeaders {
        const timestamp = Math.floor(new Date().getTime() / 1000);
        const tokenSig = `${timestamp}${this.userName}#${this.token}#${this.userName}${timestamp}`;
        const ticket = createHash("sha256").update(tokenSig).digest("hex").toUpperCase();
        return {
            "TCA-USERID": this.userName,
            "TCA-TICKET": ticket,
            "TCA-TIMESTAMP": timestamp.toString(),
        }
    }
    buildUrlWithNestedParams(url: string, params: any): string {
        const paramsObj = new URLSearchParams();
        
        // 递归处理嵌套对象
        function addParams(prefix: string, value: any) {
        if (value === undefined || value === null) {
            return;
        }
        
        if (typeof value === 'object' && !Array.isArray(value)) {
            Object.entries(value).forEach(([key, val]) => {
            addParams(`${prefix}[${key}]`, val);
            });
        } else if (Array.isArray(value)) {
            value.forEach((item, index) => {
            addParams(`${prefix}[${index}]`, item);
            });
        } else {
            paramsObj.append(prefix, String(value));
        }
        }
        
        Object.entries(params).forEach(([key, value]) => {
        addParams(key, value);
        });
        
        return paramsObj.toString() ? `${url}?${paramsObj}` : url;
    }
  

    async request<T>(path: string, method: string, data?: any, params?: any): Promise<T> {
        const norBaseUrl = this.baseUrl.endsWith("/") ? this.baseUrl : `${this.baseUrl}/`;
        const norPath = path.startsWith("/") ? path.slice(1) : path;
        const url = new URL(norPath, norBaseUrl);
        const header = this.getHeader();
        const options: RequestOptions = {
            method,
            headers: {
                "Content-Type": "application/json",
                "TCA-USERID": header['TCA-USERID'],
                "TCA-TICKET": header['TCA-TICKET'],
                "TCA-TIMESTAMP": header['TCA-TIMESTAMP']
            },
        };
        if (data) {
            options.body = JSON.stringify(data);
        }
        if (params) {
            url.search = this.buildUrlWithNestedParams(url.search, params);
        }
        try {
            const response = await fetch(url, options)
            if (!response.ok) {
                throw new Error(`Request failed with status ${response.status}, ${options.body}, ${options.headers}, ${options.method}, ${url}`);
            }
            const json = await response.json();
            return json as T;
        } catch (error) {
            throw error;
        }
    }
}
