if(!customElements.get('input-integer'))
customElements.define(
  "input-integer",
  class extends HTMLElement {
    static formAssociated = true;
    static get observedAttributes() {
      return ["required", "value", "min", "max", "name"];
    }

    _attrs = {};
    _internals;
    _defaultValue = "";

    constructor() {
      super();
      this._internals = this.attachInternals();
      this._internals.role = "textbox";
      this.tabindex = 0;
    }

    connectedCallback() {
      const shadowRoot = this.attachShadow({
        mode: "open",
        delegatesFocus: true
      });
      
      shadowRoot.innerHTML = `
        <style>
          input[type=number]::-webkit-outer-spin-button,
          input[type=number]::-webkit-inner-spin-button {
              -webkit-appearance: none;
              margin: 0;
          }

          input[type=number] {
              -moz-appearance:textfield; /* Firefox */
              font-size: 18px;
              text-align: center;
          }
          button {
            border-radius:50%;
            width: 35px;
            height: 35px;
            border: solid 1px;
            font-size: 12px;
            user-select: none;
          }
          .container {
            display: inline-block;
          }
        </style>
        <div class="container">
          <button class="decrease-button">▼</button>
          <input type="number" class="no-spinner">
          <button class="increase-button">▲</button>
        </div>`;

      this.$input = this.shadowRoot.querySelector("input");
      this.setProps();
      this._defaultValue = this.$input.value;
      this._internals.setFormValue(this.value);
      this._internals.setValidity(
        this.$input.validity,
        this.$input.validationMessage,
        this.$input
      );
      this.$input.addEventListener("input", () => this.handleInput());
      this.shadowRoot.querySelector('button.increase-button').onclick = () => this.increase();
      this.shadowRoot.querySelector('button.decrease-button').onclick = () => this.decrease();

      const form = document.querySelector('form');
      form.addEventListener('formdata', ({formData}) => {
        formData.append(this.$input.name, this.$input.value);
      });
    }

    attributeChangedCallback(name, prev, next) {
      this._attrs[name] = next;
    }

    formDisabledCallback(disabled) {
      this.$input.disabled = disabled;
    }

    formResetCallback() {
      this.$input.value = this._defaultValue;
    }

    checkValidity() {
      return this._internals.checkValidity();
    }

    reportValidity() {
      return this._internals.reportValidity();
    }

    get validity() {
      return this._internals.validity;
    }

    get validationMessage() {
      return this._internals.validationMessage;
    }

    setProps() {
      // prevent any errors in case the input isn't set
      if (!this.$input) {
        return;
      }

      for (let prop in this._attrs) {
        switch (prop) {
          case "name":
            this.$input.name = this._attrs[prop];
            break;
          case "min":
            this.$input.setAttribute("min",this._attrs[prop]);
            break;
          case "max":
            this.$input.setAttribute("max",this._attrs[prop]);
            break;
          case "value":
            this.$input.value = this._attrs[prop];
            break;
          case "required":
            const required = this._attrs[prop];
            this.$input.toggleAttribute(
              "required",
              required === "true" || required === ""
            );
            break;
        }
      }

      this._attrs = {};
    }

    handleInput() {
      this._internals.setValidity(
        this.$input.validity,
        this.$input.validationMessage,
        this.$input
      );
      this._internals.setFormValue(this.value);
    }

  currentValue() { 
    return  (parseInt(this.$input.value) || 0);
  };

  increase() {
    let currentValue = this.currentValue(0);
    if (currentValue + 1 <= this.$input.getAttribute("max")) {
      this.$input.value = currentValue + 1;
    }
    else {
      this.$input.value = this.$input.getAttribute("max");
    }
    this.handleInput();
  }

  decrease() {
    let currentValue = this.currentValue();
    if (this.$input.value == "") {
      this.$input.value = this.$input.getAttribute("max");
    }
    else if (currentValue - 1 >= this.$input.getAttribute("min")) {
      this.$input.value = currentValue - 1;
    }
    else {
      this.$input.value = this.$input.getAttribute("min");
    }
    this.handleInput();
  }
});
