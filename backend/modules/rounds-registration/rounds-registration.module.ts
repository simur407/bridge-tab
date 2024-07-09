import { Module } from "@packages/modules";
import { HttpsFunction } from "firebase-functions/v2/https";

export class RoundsRegistrationModule implements Module {
  name = "rounds-registration";
  controllers(): Record<string, HttpsFunction> {
    throw new Error("Method not implemented.");
  }
}