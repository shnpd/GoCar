import * as $protobuf from "protobufjs";

// Common aliases
const $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

export const coolcar = $root.coolcar = (() => {

    /**
     * Namespace coolcar.
     * @exports coolcar
     * @namespace
     */
    const coolcar = {};

    coolcar.Location = (function() {

        /**
         * Properties of a Location.
         * @memberof coolcar
         * @interface ILocation
         * @property {number|null} [latitude] Location latitude
         * @property {number|null} [longitude] Location longitude
         */

        /**
         * Constructs a new Location.
         * @memberof coolcar
         * @classdesc Represents a Location.
         * @implements ILocation
         * @constructor
         * @param {coolcar.ILocation=} [properties] Properties to set
         */
        function Location(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Location latitude.
         * @member {number} latitude
         * @memberof coolcar.Location
         * @instance
         */
        Location.prototype.latitude = 0;

        /**
         * Location longitude.
         * @member {number} longitude
         * @memberof coolcar.Location
         * @instance
         */
        Location.prototype.longitude = 0;

        /**
         * Encodes the specified Location message. Does not implicitly {@link coolcar.Location.verify|verify} messages.
         * @function encode
         * @memberof coolcar.Location
         * @static
         * @param {coolcar.ILocation} message Location message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Location.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.latitude != null && Object.hasOwnProperty.call(message, "latitude"))
                writer.uint32(/* id 1, wireType 1 =*/9).double(message.latitude);
            if (message.longitude != null && Object.hasOwnProperty.call(message, "longitude"))
                writer.uint32(/* id 2, wireType 1 =*/17).double(message.longitude);
            return writer;
        };

        /**
         * Creates a Location message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof coolcar.Location
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {coolcar.Location} Location
         */
        Location.fromObject = function fromObject(object) {
            if (object instanceof $root.coolcar.Location)
                return object;
            let message = new $root.coolcar.Location();
            if (object.latitude != null)
                message.latitude = Number(object.latitude);
            if (object.longitude != null)
                message.longitude = Number(object.longitude);
            return message;
        };

        /**
         * Creates a plain object from a Location message. Also converts values to other types if specified.
         * @function toObject
         * @memberof coolcar.Location
         * @static
         * @param {coolcar.Location} message Location
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Location.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.latitude = 0;
                object.longitude = 0;
            }
            if (message.latitude != null && message.hasOwnProperty("latitude"))
                object.latitude = options.json && !isFinite(message.latitude) ? String(message.latitude) : message.latitude;
            if (message.longitude != null && message.hasOwnProperty("longitude"))
                object.longitude = options.json && !isFinite(message.longitude) ? String(message.longitude) : message.longitude;
            return object;
        };

        /**
         * Converts this Location to JSON.
         * @function toJSON
         * @memberof coolcar.Location
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Location.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return Location;
    })();

    /**
     * TripStatus enum.
     * @name coolcar.TripStatus
     * @enum {number}
     * @property {number} TS_NOT_SPECIFIED=0 TS_NOT_SPECIFIED value
     * @property {number} NOT_STARTED=1 NOT_STARTED value
     * @property {number} IN_PROGRESS=2 IN_PROGRESS value
     * @property {number} FINISHED=3 FINISHED value
     * @property {number} PAID=4 PAID value
     */
    coolcar.TripStatus = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "TS_NOT_SPECIFIED"] = 0;
        values[valuesById[1] = "NOT_STARTED"] = 1;
        values[valuesById[2] = "IN_PROGRESS"] = 2;
        values[valuesById[3] = "FINISHED"] = 3;
        values[valuesById[4] = "PAID"] = 4;
        return values;
    })();

    coolcar.Trip = (function() {

        /**
         * Properties of a Trip.
         * @memberof coolcar
         * @interface ITrip
         * @property {string|null} [start] Trip start
         * @property {coolcar.ILocation|null} [startPos] Trip startPos
         * @property {Array.<coolcar.ILocation>|null} [pathLocations] Trip pathLocations
         * @property {string|null} [end] Trip end
         * @property {coolcar.ILocation|null} [endPos] Trip endPos
         * @property {number|null} [durationSec] Trip durationSec
         * @property {number|null} [feeCent] Trip feeCent
         * @property {coolcar.TripStatus|null} [status] Trip status
         */

        /**
         * Constructs a new Trip.
         * @memberof coolcar
         * @classdesc Represents a Trip.
         * @implements ITrip
         * @constructor
         * @param {coolcar.ITrip=} [properties] Properties to set
         */
        function Trip(properties) {
            this.pathLocations = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Trip start.
         * @member {string} start
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.start = "";

        /**
         * Trip startPos.
         * @member {coolcar.ILocation|null|undefined} startPos
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.startPos = null;

        /**
         * Trip pathLocations.
         * @member {Array.<coolcar.ILocation>} pathLocations
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.pathLocations = $util.emptyArray;

        /**
         * Trip end.
         * @member {string} end
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.end = "";

        /**
         * Trip endPos.
         * @member {coolcar.ILocation|null|undefined} endPos
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.endPos = null;

        /**
         * Trip durationSec.
         * @member {number} durationSec
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.durationSec = 0;

        /**
         * Trip feeCent.
         * @member {number} feeCent
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.feeCent = 0;

        /**
         * Trip status.
         * @member {coolcar.TripStatus} status
         * @memberof coolcar.Trip
         * @instance
         */
        Trip.prototype.status = 0;

        /**
         * Encodes the specified Trip message. Does not implicitly {@link coolcar.Trip.verify|verify} messages.
         * @function encode
         * @memberof coolcar.Trip
         * @static
         * @param {coolcar.ITrip} message Trip message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Trip.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.start != null && Object.hasOwnProperty.call(message, "start"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.start);
            if (message.end != null && Object.hasOwnProperty.call(message, "end"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.end);
            if (message.durationSec != null && Object.hasOwnProperty.call(message, "durationSec"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.durationSec);
            if (message.feeCent != null && Object.hasOwnProperty.call(message, "feeCent"))
                writer.uint32(/* id 4, wireType 0 =*/32).int32(message.feeCent);
            if (message.startPos != null && Object.hasOwnProperty.call(message, "startPos"))
                $root.coolcar.Location.encode(message.startPos, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
            if (message.endPos != null && Object.hasOwnProperty.call(message, "endPos"))
                $root.coolcar.Location.encode(message.endPos, writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
            if (message.pathLocations != null && message.pathLocations.length)
                for (let i = 0; i < message.pathLocations.length; ++i)
                    $root.coolcar.Location.encode(message.pathLocations[i], writer.uint32(/* id 7, wireType 2 =*/58).fork()).ldelim();
            if (message.status != null && Object.hasOwnProperty.call(message, "status"))
                writer.uint32(/* id 8, wireType 0 =*/64).int32(message.status);
            return writer;
        };

        /**
         * Creates a Trip message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof coolcar.Trip
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {coolcar.Trip} Trip
         */
        Trip.fromObject = function fromObject(object) {
            if (object instanceof $root.coolcar.Trip)
                return object;
            let message = new $root.coolcar.Trip();
            if (object.start != null)
                message.start = String(object.start);
            if (object.startPos != null) {
                if (typeof object.startPos !== "object")
                    throw TypeError(".coolcar.Trip.startPos: object expected");
                message.startPos = $root.coolcar.Location.fromObject(object.startPos);
            }
            if (object.pathLocations) {
                if (!Array.isArray(object.pathLocations))
                    throw TypeError(".coolcar.Trip.pathLocations: array expected");
                message.pathLocations = [];
                for (let i = 0; i < object.pathLocations.length; ++i) {
                    if (typeof object.pathLocations[i] !== "object")
                        throw TypeError(".coolcar.Trip.pathLocations: object expected");
                    message.pathLocations[i] = $root.coolcar.Location.fromObject(object.pathLocations[i]);
                }
            }
            if (object.end != null)
                message.end = String(object.end);
            if (object.endPos != null) {
                if (typeof object.endPos !== "object")
                    throw TypeError(".coolcar.Trip.endPos: object expected");
                message.endPos = $root.coolcar.Location.fromObject(object.endPos);
            }
            if (object.durationSec != null)
                message.durationSec = object.durationSec | 0;
            if (object.feeCent != null)
                message.feeCent = object.feeCent | 0;
            switch (object.status) {
            case "TS_NOT_SPECIFIED":
            case 0:
                message.status = 0;
                break;
            case "NOT_STARTED":
            case 1:
                message.status = 1;
                break;
            case "IN_PROGRESS":
            case 2:
                message.status = 2;
                break;
            case "FINISHED":
            case 3:
                message.status = 3;
                break;
            case "PAID":
            case 4:
                message.status = 4;
                break;
            }
            return message;
        };

        /**
         * Creates a plain object from a Trip message. Also converts values to other types if specified.
         * @function toObject
         * @memberof coolcar.Trip
         * @static
         * @param {coolcar.Trip} message Trip
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Trip.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.pathLocations = [];
            if (options.defaults) {
                object.start = "";
                object.end = "";
                object.durationSec = 0;
                object.feeCent = 0;
                object.startPos = null;
                object.endPos = null;
                object.status = options.enums === String ? "TS_NOT_SPECIFIED" : 0;
            }
            if (message.start != null && message.hasOwnProperty("start"))
                object.start = message.start;
            if (message.end != null && message.hasOwnProperty("end"))
                object.end = message.end;
            if (message.durationSec != null && message.hasOwnProperty("durationSec"))
                object.durationSec = message.durationSec;
            if (message.feeCent != null && message.hasOwnProperty("feeCent"))
                object.feeCent = message.feeCent;
            if (message.startPos != null && message.hasOwnProperty("startPos"))
                object.startPos = $root.coolcar.Location.toObject(message.startPos, options);
            if (message.endPos != null && message.hasOwnProperty("endPos"))
                object.endPos = $root.coolcar.Location.toObject(message.endPos, options);
            if (message.pathLocations && message.pathLocations.length) {
                object.pathLocations = [];
                for (let j = 0; j < message.pathLocations.length; ++j)
                    object.pathLocations[j] = $root.coolcar.Location.toObject(message.pathLocations[j], options);
            }
            if (message.status != null && message.hasOwnProperty("status"))
                object.status = options.enums === String ? $root.coolcar.TripStatus[message.status] : message.status;
            return object;
        };

        /**
         * Converts this Trip to JSON.
         * @function toJSON
         * @memberof coolcar.Trip
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Trip.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return Trip;
    })();

    coolcar.GetTripRequest = (function() {

        /**
         * Properties of a GetTripRequest.
         * @memberof coolcar
         * @interface IGetTripRequest
         * @property {string|null} [id] GetTripRequest id
         */

        /**
         * Constructs a new GetTripRequest.
         * @memberof coolcar
         * @classdesc Represents a GetTripRequest.
         * @implements IGetTripRequest
         * @constructor
         * @param {coolcar.IGetTripRequest=} [properties] Properties to set
         */
        function GetTripRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * GetTripRequest id.
         * @member {string} id
         * @memberof coolcar.GetTripRequest
         * @instance
         */
        GetTripRequest.prototype.id = "";

        /**
         * Encodes the specified GetTripRequest message. Does not implicitly {@link coolcar.GetTripRequest.verify|verify} messages.
         * @function encode
         * @memberof coolcar.GetTripRequest
         * @static
         * @param {coolcar.IGetTripRequest} message GetTripRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetTripRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            return writer;
        };

        /**
         * Creates a GetTripRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof coolcar.GetTripRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {coolcar.GetTripRequest} GetTripRequest
         */
        GetTripRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.coolcar.GetTripRequest)
                return object;
            let message = new $root.coolcar.GetTripRequest();
            if (object.id != null)
                message.id = String(object.id);
            return message;
        };

        /**
         * Creates a plain object from a GetTripRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof coolcar.GetTripRequest
         * @static
         * @param {coolcar.GetTripRequest} message GetTripRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetTripRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.id = "";
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            return object;
        };

        /**
         * Converts this GetTripRequest to JSON.
         * @function toJSON
         * @memberof coolcar.GetTripRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetTripRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return GetTripRequest;
    })();

    coolcar.GetTripResponse = (function() {

        /**
         * Properties of a GetTripResponse.
         * @memberof coolcar
         * @interface IGetTripResponse
         * @property {string|null} [id] GetTripResponse id
         * @property {coolcar.ITrip|null} [trip] GetTripResponse trip
         */

        /**
         * Constructs a new GetTripResponse.
         * @memberof coolcar
         * @classdesc Represents a GetTripResponse.
         * @implements IGetTripResponse
         * @constructor
         * @param {coolcar.IGetTripResponse=} [properties] Properties to set
         */
        function GetTripResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * GetTripResponse id.
         * @member {string} id
         * @memberof coolcar.GetTripResponse
         * @instance
         */
        GetTripResponse.prototype.id = "";

        /**
         * GetTripResponse trip.
         * @member {coolcar.ITrip|null|undefined} trip
         * @memberof coolcar.GetTripResponse
         * @instance
         */
        GetTripResponse.prototype.trip = null;

        /**
         * Encodes the specified GetTripResponse message. Does not implicitly {@link coolcar.GetTripResponse.verify|verify} messages.
         * @function encode
         * @memberof coolcar.GetTripResponse
         * @static
         * @param {coolcar.IGetTripResponse} message GetTripResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetTripResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.trip != null && Object.hasOwnProperty.call(message, "trip"))
                $root.coolcar.Trip.encode(message.trip, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
            return writer;
        };

        /**
         * Creates a GetTripResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof coolcar.GetTripResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {coolcar.GetTripResponse} GetTripResponse
         */
        GetTripResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.coolcar.GetTripResponse)
                return object;
            let message = new $root.coolcar.GetTripResponse();
            if (object.id != null)
                message.id = String(object.id);
            if (object.trip != null) {
                if (typeof object.trip !== "object")
                    throw TypeError(".coolcar.GetTripResponse.trip: object expected");
                message.trip = $root.coolcar.Trip.fromObject(object.trip);
            }
            return message;
        };

        /**
         * Creates a plain object from a GetTripResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof coolcar.GetTripResponse
         * @static
         * @param {coolcar.GetTripResponse} message GetTripResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetTripResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.trip = null;
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.trip != null && message.hasOwnProperty("trip"))
                object.trip = $root.coolcar.Trip.toObject(message.trip, options);
            return object;
        };

        /**
         * Converts this GetTripResponse to JSON.
         * @function toJSON
         * @memberof coolcar.GetTripResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetTripResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return GetTripResponse;
    })();

    coolcar.TripService = (function() {

        /**
         * Constructs a new TripService service.
         * @memberof coolcar
         * @classdesc Represents a TripService
         * @extends $protobuf.rpc.Service
         * @constructor
         * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
         * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
         * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
         */
        function TripService(rpcImpl, requestDelimited, responseDelimited) {
            $protobuf.rpc.Service.call(this, rpcImpl, requestDelimited, responseDelimited);
        }

        (TripService.prototype = Object.create($protobuf.rpc.Service.prototype)).constructor = TripService;

        /**
         * Callback as used by {@link coolcar.TripService#getTrip}.
         * @memberof coolcar.TripService
         * @typedef GetTripCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {coolcar.GetTripResponse} [response] GetTripResponse
         */

        /**
         * Calls GetTrip.
         * @function getTrip
         * @memberof coolcar.TripService
         * @instance
         * @param {coolcar.IGetTripRequest} request GetTripRequest message or plain object
         * @param {coolcar.TripService.GetTripCallback} callback Node-style callback called with the error, if any, and GetTripResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(TripService.prototype.getTrip = function getTrip(request, callback) {
            return this.rpcCall(getTrip, $root.coolcar.GetTripRequest, $root.coolcar.GetTripResponse, request, callback);
        }, "name", { value: "GetTrip" });

        /**
         * Calls GetTrip.
         * @function getTrip
         * @memberof coolcar.TripService
         * @instance
         * @param {coolcar.IGetTripRequest} request GetTripRequest message or plain object
         * @returns {Promise<coolcar.GetTripResponse>} Promise
         * @variation 2
         */

        return TripService;
    })();

    return coolcar;
})();