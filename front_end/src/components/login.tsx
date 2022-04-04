import React from "react";
import TextField from '@mui/material/TextField';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import { red, grey } from '@mui/material/colors';

export const LoginScreen = () => {
    return (
        <Card sx={{ display: 'flex', height: '100vh' }}>
            <CardContent sx={{bgcolor: red[500]}}>
                <Typography sx={{ fontSize: 20 }} color="text.secondary" gutterBottom>
                    欢迎来到欢乐升级!
                </Typography>
            </CardContent>
            <CardContent sx={{ flex: '1 0 auto', bgcolor: grey[500] }}>
                <TextField 
                    label="请输入一个昵称开始游戏" 
                    color="error" 
                    focused 
                    margin='normal'
                    sx={{
                        fontSize: "40"
                    }}
                />
            </CardContent>
        </Card>
    )
}
