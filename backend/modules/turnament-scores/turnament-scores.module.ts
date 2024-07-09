import { Module } from "@packages/modules";
import { HttpsFunction } from "firebase-functions/v2/https";

export class TurnamentScoresModule implements Module {
  name = "turnament-scores";
  controllers(): Record<string, HttpsFunction> {
    throw new Error("Method not implemented.");
  }
}