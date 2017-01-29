local krpc = require 'krpc'
local conn = krpc:connect('User Interface Example')
local canvas = conn.ui.stock_canvas

-- Get the size of the game window in pixels
local screen_size = canvas.rect_transform.size

-- Add a panel to contain the UI elements
local panel = canvas:add_panel()

-- Position the panel on the left of the screen
local rect = panel.rect_transform
rect.size = {200, 100}
rect.position = {110-(screen_size[1]/2), 0}

-- Add a button to set the throttle to maximum
local button = panel:add_button("Full Throttle")
button.rect_transform.position = {0, 20}

-- Add some text displaying the total engine thrust
local text = panel:add_text("Thrust: 0 kN")
text.rect_transform.position = {0, -20}
text.color = {1, 1, 1}
text.size = 18

local vessel = conn.space_center.active_vessel
while True do
    -- Handle the throttle button being clicked
    if button.clicked then
        vessel.control.throttle = 1
        button.clicked = False
    end

    -- Update the thrust text
    text.content = 'Thrust: ' .. (vessel.thrust/1000) .. ' kN'

    time.sleep(0.1)
end