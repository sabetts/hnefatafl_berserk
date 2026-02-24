import random
from PIL import Image

tiles = [Image.open(f'tile_{i:02}.png') for i in range(30)]

corner_bot_left = Image.open('wood_corner_bot_left.png')
corner_bot_right = corner_bot_left.transpose(Image.Transpose.FLIP_LEFT_RIGHT)
corner_top_left = Image.open('wood_corner_top_left.png')
corner_top_right = corner_top_left.transpose(Image.Transpose.FLIP_LEFT_RIGHT)

edge_bottom = Image.open('wood_edge_bottom.png')
edge_top = Image.open('wood_edge_top.png')
edge_right = Image.open('wood_edge_right.png')
edge_left = Image.open('wood_edge_left.png')

edge_top2 = edge_bottom.transpose(Image.Transpose.FLIP_TOP_BOTTOM)
edge_bottom2 = edge_top.transpose(Image.Transpose.FLIP_TOP_BOTTOM)
edge_left2 = edge_right.transpose(Image.Transpose.FLIP_LEFT_RIGHT)
edge_right2 = edge_left.transpose(Image.Transpose.FLIP_LEFT_RIGHT)

ofsx = 138
ofsy = 138

print(tiles[0].width, tiles[0].height)

# width = tiles[0].width * 11 + corner_bot_left.width + corner_bot_right.width
# height = tiles[0].height * 11 + corner_bot_left.height + corner_bot_right.height
width = tiles[0].width * 11 + ofsx + ofsx
height = tiles[0].height * 11 + ofsy +ofsy-15

bg = (50,50,50)
output = Image.new("RGB", (width,height), bg)


for x in range(0,11):
    for y in range(0,11):
        #i = random.randint(0, len(tiles)-1)
        i = (y*11+x)%len(tiles)
        t = tiles[i]
        output.paste(t, (ofsx+t.width*x, ofsy+t.width*y))

output.paste(edge_top, (corner_bot_left.width, 0))
output.paste(edge_top2, (corner_bot_left.width + edge_top.width, 3))
output.paste(edge_top, (corner_bot_left.width + edge_top.width+edge_top2.width, 0))
output.paste(edge_top2, (corner_bot_left.width + edge_top.width+edge_top2.width*2, 5))

output.paste(edge_bottom, (corner_bot_left.width, output.height-edge_bottom.height))
output.paste(edge_bottom2, (corner_bot_left.width + edge_bottom.width, output.height-edge_bottom2.height+3))
output.paste(edge_bottom, (corner_bot_left.width + edge_bottom.width+edge_bottom2.width, output.height-edge_bottom.height))
output.paste(edge_bottom2, (corner_bot_left.width + edge_bottom.width+edge_bottom2.width*2, output.height-edge_bottom2.height))


output.paste(edge_left, (0, corner_bot_left.height))
output.paste(edge_left2, (2, corner_bot_left.height + edge_left.height))
output.paste(edge_left, (0, corner_bot_left.height + edge_left.height+edge_left2.height))
output.paste(edge_left2, (2, corner_bot_left.height + edge_left.height+edge_left2.height*2))

output.paste(edge_right, (output.width-edge_right.width+2, corner_bot_right.height))
output.paste(edge_right2, (output.width-edge_right.width, corner_bot_right.height + edge_right.height))
output.paste(edge_right, (output.width-edge_right.width, corner_bot_right.height + edge_right.height+edge_right2.height))
output.paste(edge_right2, (output.width-edge_right.width, corner_bot_right.height + edge_right.height+edge_right2.height*2))




output.paste(corner_bot_left, (0,height-corner_bot_left.height), corner_bot_left)
output.paste(corner_bot_right, (width-corner_bot_right.height,height-corner_bot_right.height), corner_bot_right)
output.paste(corner_top_left, (0,0), corner_top_left)
output.paste(corner_top_right, (width-corner_bot_right.height,0), corner_top_right)


# mockup a board position
if False:
    throne = Image.open('throne.png')
    attacker = Image.open('attacker.png')
    commander = Image.open('commander.png')
    defender = Image.open('defender.png')
    king = Image.open('shortking2.png')
    knight = Image.open('knight.png')

    def place(img, tx, ty):
        cx = (tiles[0].width-img.width)//2
        cy = (tiles[0].height-img.height)//2
        output.paste(img, (ofsx + cx + tx * tiles[0].width, ofsy + cy + ty * tiles[0].height), img)

    for i in range(3,8):
        place(attacker, i, 0)
        place(attacker, i, 10)
        place(attacker, 0, i)
        place(attacker, 10, i)

    place(commander, 1, 5)
    place(commander, 9, 5)
    place(commander, 5, 1)
    place(commander, 5, 9)

    place(defender, 3, 5)
    place(defender, 7, 5)
    place(defender, 5, 3)
    place(defender, 5, 7)
    place(defender, 5, 4)
    place(defender, 5, 6)

    place(defender, 4, 4)
    place(defender, 4, 5)
    place(defender, 4, 6)
    place(knight, 6, 4)
    place(defender, 6, 5)
    place(defender, 6, 6)


    output.paste(throne, (ofsx+t.width*5+(t.width-throne.width)//2, ofsy+t.height*6-throne.height), throne)
    place(king, 5,5)


output.save('board_final.png')
