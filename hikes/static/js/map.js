import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

const map = new maplibregl.Map({
    container: "map",
    style: "map-style.json",
    center: [-96.5, 39.8], // [lng, lat]
    zoom: 3.3,
    attributionControl: false // 👈 disables attribution box
});

// Test hard coded geo locations
const locations = [
    // {name: "New York City", coords: [-74.006, 40.7128]},
    // {name: "Los Angeles", coords: [-118.2437, 34.0522]},
    // {name: "Chicago", coords: [-87.6298, 41.8781]}
];

function markLocation(loc) {
    const el = document.createElement("div");
    el.className = "marker";
    el.style.cssText = `
    background: red;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    border: 2px solid white;
  `;

    new maplibregl.Marker(el)
        .setLngLat(loc.coords)
        .setPopup(new maplibregl.Popup({ offset: 25 }).setText(loc.name))
        .addTo(map);

}

// Setup scss
const root = document.documentElement;
const backgroundColor = getComputedStyle(root).getPropertyValue('--color-background').trim();

console.log("Before map loaded!");

map.on("load", () => {

    locations.forEach(loc => {
        markLocation(loc);
    });

    // iframe map - get message from theme-switch.js
    window.addEventListener("message", (event) => {
        console.log("Map received:", event.data);
        if (event.data === "theme-light") {
            map.setPaintProperty('background', 'background-color', '#f5f5f5');
            map.setPaintProperty('state-fills', 'fill-color', '#dddddd');
            map.setPaintProperty('state-borders', 'line-color', '#1f1f1f');
        } else if (event.data === "theme-dark") {
            map.setPaintProperty('background', 'background-color', '#333333');
            map.setPaintProperty('state-fills', 'fill-color', '#525252');
            map.setPaintProperty('state-borders', 'line-color', '#ededed');
        }
    });

    // Tell the calling page that the map is ready
    window.parent.postMessage({ type: "map-ready" }, "*");

    // For the map without iframe
    const dark_mode_btn = document.getElementById("dark_mode_btn");
    const light_mode_btn = document.getElementById("light_mode_btn");
    if (dark_mode_btn) {
        console.log("Dark mode button found");
        dark_mode_btn.addEventListener('click', function () {
            console.log("[map.js] Dark mode");
            map.setPaintProperty('background', 'background-color', '#333333');
        });
    }
    if (light_mode_btn) {
        console.log("Light mode button found");
        light_mode_btn.addEventListener('click', function () {
            console.log("[map.js] Light mode");
            map.setPaintProperty('background', 'background-color', '#f7f7f7');
        });
    }

    // mark (hike) locations on the map
    fetch('locations.json')
        .then(response => response.json())
        .then(data => {
            data.forEach(loc => {
                markLocation(loc);
            });
        })
        .catch(err => console.error('Error loading JSON:', err));

    // fit map based on continental US
    const continentalUSBounds = [
        [-125.0, 24.5],  // Southwest corner (lon, lat)
        [-66.9, 49.5]    // Northeast corner (lon, lat)
    ];
    map.fitBounds(continentalUSBounds, {
        padding: 40,
        animate: true,
        duration: 1000
    });
    // Re-fit on window resize
    window.addEventListener('resize', () => {
        map.fitBounds(continentalUSBounds, {
            padding: 40,
            animate: true,
            duration: 1000
        });
    });
});
