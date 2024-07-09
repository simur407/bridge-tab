export type TurnamentState = {
  id: string;
  name: string;
};

export class TurnamentCreated {
  readonly id: string;
  readonly name: string;

  constructor(event: TurnamentCreated) {
    this.id = event.id;
    this.name = event.name;
  }
}

type TurnamentEvents = TurnamentCreated;

export class Turnament {
  private id: string;

  constructor(readonly state: TurnamentState) {
    this.id = state.id;
  }

  static create(state: TurnamentState) {
    const entity = new Turnament(state);
    entity.events.push(new TurnamentCreated(state));
    return entity;
  }

  // Domain events
  private events: TurnamentEvents[] = [];
  getEvents() {
    return this.events;
  }
}

export interface TurnamentRepository {
  load(id: string): Promise<Turnament | undefined>;
  save(turnament: Turnament): Promise<void>;
}

export class InMemoryTurnamentRepository implements TurnamentRepository {
  private turnaments: Turnament[] = [];
  async load(id: string): Promise<Turnament | undefined> {
    const turnament = this.turnaments.find((turnament) => turnament.state.id === id);
    if (!turnament) {
      return undefined;
    }
    return turnament;
  }
  async save(turnament: Turnament): Promise<void> {
    this.turnaments = this.turnaments.filter((t) => t.state.id !== turnament.state.id);
    this.turnaments.push(turnament);
  }
}
