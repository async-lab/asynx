import webapi from "./gocliRequest"
import * as components from "./custodianComponents"
export * from "./custodianComponents"

/**
 * @description 
 * @param params
 */
export function hello(params: components.RequestParams, name: string) {
	return webapi.get<components.Response>(`/from/${name}`, params)
}
