export type CallbackID = string;

export interface SendOptions {
  model?: string;
  body?: string;
}

export interface ClientPlugin {
  send(
    options: SendOptions,
    callback: ClientCallback
  ): Promise<CallbackID>;
}


export type ClientCallback = (
  (response: string, err?: any) => boolean
)

