export type CallbackID = string;

export interface ClientPlugin {
  send(
    url: string,
    body: string,
    callback: ClientCallback): Promise<CallbackID>;
}


export type ClientCallback = (
  (response: string, err?: any) => boolean
)

