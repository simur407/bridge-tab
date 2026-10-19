/**
 * Custom element <input-integer> that provides an integer input field with increase and decrease buttons.
 */
if(!customElements.get('input-integer'))
customElements.define(
  "input-integer",
  class extends HTMLElement {
    static formAssociated = true;
    /**
     * Gets the list of attributes to observe for changes.
     * @returns {string[]} Array of attribute names to observe
     */
    static get observedAttributes() {
      return ["required", "value", "min", "max", "name"];
    }

    _attrs = {};
    _internals;
    _defaultValue = "";

    /**
     * Constructs a new input-integer custom element.
     * Initializes internals and sets up the element as a form-associated custom element.
     */
    constructor() {
      super();
      this._internals = this.attachInternals();
      this._internals.role = "textbox";
      this.tabindex = 0;
    }

    /**
     * Lifecycle callback invoked when the element is connected to the DOM.
     * Creates the shadow DOM, sets up event listeners, and initializes the form element.
     */
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

    /**
     * Lifecycle callback invoked when an observed attribute changes.
     * @param {string} name - The name of the attribute that changed
     * @param {string} prev - The previous value of the attribute
     * @param {string} next - The new value of the attribute
     */
    attributeChangedCallback(name, prev, next) {
      this._attrs[name] = next;
    }

    /**
     * Lifecycle callback invoked when the form's disabled state changes.
     * @param {boolean} disabled - Whether the form is disabled
     */
    formDisabledCallback(disabled) {
      this.$input.disabled = disabled;
    }

    /**
     * Lifecycle callback invoked when the form is reset.
     * Restores the input to its default value.
     */
    formResetCallback() {
      this.$input.value = this._defaultValue;
    }

    /**
     * Checks whether the element meets its validation constraints.
     * @returns {boolean} True if the element is valid, false otherwise
     */
    checkValidity() {
      return this._internals.checkValidity();
    }

    /**
     * Checks validity and reports validation errors to the user.
     * @returns {boolean} True if the element is valid, false otherwise
     */
    reportValidity() {
      return this._internals.reportValidity();
    }

    /**
     * Gets the validity state of the element.
     * @returns {ValidityState} The validity state object
     */
    get validity() {
      return this._internals.validity;
    }

    /**
     * Gets the validation message for the element.
     * @returns {string} The validation message
     */
    get validationMessage() {
      return this._internals.validationMessage;
    }

    /**
     * Applies the cached attributes to the input element.
     * Handles name, min, max, value, and required attributes.
     */
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

    /**
     * Handles input changes by updating the element's validity state and form value.
     */
    handleInput() {
      this._internals.setValidity(
        this.$input.validity,
        this.$input.validationMessage,
        this.$input
      );
      this._internals.setFormValue(this.value);
    }

  /**
   * Gets the current integer value of the input.
   * @returns {number} The parsed integer value, or 0 if the value is not a valid number
   */
  currentValue() { 
    return  (parseInt(this.$input.value) || 0);
  };

  /**
   * Increases the input value by 1, respecting the maximum constraint.
   * If the value is already at maximum, it remains unchanged.
   */
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

  /**
   * Decreases the input value by 1, respecting the minimum constraint.
   * If the value is empty, sets it to the maximum value.
   * If the value is already at minimum, it remains unchanged.
   */
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
