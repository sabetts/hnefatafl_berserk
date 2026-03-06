from PIL import Image

board = Image.open('hnefetafle_board_mockup.png')

ofs_x = 141
ofs_y = 143

tile_width = 106
tile_height = 107


coords = [
    (2,0), (3,0), (4,0),
    (2,1), (3,1), (4,1),
    (0,2), (1,2), (2,2), (4,2), (5,2), (6,2),
    (0,3), (1,3), (2,3), (4,3), (5,3), (6,3),
    (0,4), (1,4), (2,4), (4,4), (5,4), (6,4),
    (2,5), (3,5), (4,5),
    (2,6), (3,6), (4,6),
]

c = 0
for tx,ty in coords:
    x = ofs_x + tx * tile_width
    y = ofs_y + ty * tile_height
    tile = board.crop((x,y,x+tile_width,y+tile_height))
    tile.save(f'tile_{c:02}.png')
    c += 1

top_wood = board.crop((353, 0, 671, 135))
bot_wood = board.crop((353, 894, 671, 1023))
lef_wood = board.crop((0, 356, 137, 677))
rig_wood = board.crop((890, 356, 1023, 677))

top_wood.save('wood_edge_top.png')
bot_wood.save('wood_edge_bottom.png')
lef_wood.save('wood_edge_left.png')
rig_wood.save('wood_edge_right.png')
