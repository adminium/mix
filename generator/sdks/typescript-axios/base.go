package typescript_axios

const baseTpl = `
/* tslint:disable */
/* eslint-disable */
/**
 *
 * No description provided (generated by Mix Openapi Generator https://github.com/gozelle/mix)
 *
 * The version of the OpenAPI document: 3.0.3
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://mix.gozelle.io).
 * Do not edit the class manually.
 */

// Some imports not used depending on template conditions
import {AxiosError, AxiosInstance, AxiosResponse, AxiosRequestConfig} from 'axios';

export class BaseAPI {
    public client: AxiosInstance;

    constructor(instance: AxiosInstance) {
        this.client = instance;
    }
}

export const requestInterceptorExample: any = [function (request: AxiosRequestConfig) {
    // 在发送请求之前做些什么
    if (request.headers) {
        request.headers['Authorization'] = ` + "`Bearer ${localStorage.getItem('access_token')}`" + `;
    }
    return request;
}, function (error: AxiosError) {
    // 对请求错误做些什么
    return Promise.reject(error);
}]

export const responseInterceptorExample: any = [function (response: AxiosResponse) {
    // 2xx 范围内的状态码都会触发该函数。
    // 对响应数据做点什么
    return response.data.result;
}, function (error: AxiosError) {
    // 超出 2xx 范围的状态码都会触发该函数。
    // 对响应错误做点什么
    if (error.response) {
        const data = error.response.data as any;
        if (error.response.status >= 400 && error.response.status < 500) {
            console.warn({ content: data ? data.message : 'custom error'})
        } else {
            console.error({ content: data ? data.message : 'system error' })
        }
    } else {
        console.error('unknown error:', JSON.stringify(error))
    }
	// 会阻止后续的业务处理逻辑
    return Promise.reject(error));
}]
`