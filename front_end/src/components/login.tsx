import React from "react";
import TextField from '@mui/material/TextField';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import { red, grey } from '@mui/material/colors';
import Button from '@mui/material/Button';

export const LoginScreen = () => {
    return (
        <Card sx={{ height: '100vh' }}>
            <CardContent sx={{bgcolor: red[500], height: '5vh'}}>
                <Typography sx={{ fontSize: 20 }} color="text.secondary" gutterBottom>
                    欢迎来到欢乐升级!
                </Typography>
            </CardContent>
            <CardContent 
                sx={{ 
                    height: '95vh',
                    bgcolor: grey[500], 
                    textAlign: "center",
                }}
            >
                <TextField 
                    label="请输入一个昵称" 
                    color="error" 
                    focused 
                    margin='normal'
                    sx={{
                        fontSize: "40"
                    }}
                />
                <CardContent
                    sx={{ 
                        bgcolor: grey[500], 
                        textAlign: "center",
                    }}
                >
                    <Button variant="contained" color="success" size="large">开始游戏</Button>
                </CardContent>
            </CardContent>
        </Card>
    )
}
